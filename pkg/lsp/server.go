package lsp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
	lsp "pkg.nimblebun.works/go-lsp"
	"src.elv.sh/pkg/diag"
	"src.elv.sh/pkg/edit/complete"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/mods/doc"
	"src.elv.sh/pkg/parse"
	"src.elv.sh/pkg/parse/cmpd"
	"src.elv.sh/pkg/parse/np"
)

var (
	errMethodNotFound = &jsonrpc2.Error{
		Code: jsonrpc2.CodeMethodNotFound, Message: "method not found",
	}
	errInvalidParams = &jsonrpc2.Error{
		Code: jsonrpc2.CodeInvalidParams, Message: "invalid params",
	}
)

type server struct {
	evaler    *eval.Evaler
	documents map[lsp.DocumentURI]document
}

type document struct {
	code      string
	parseTree parse.Tree
	parseErr  error
}

// isVariableShadowed checks if a variable is locally defined (shadowing global/builtin)
func isVariableShadowed(parseTree parse.Tree, pos int, varName string) bool {
	// Find the path to the current position
	path := np.Find(parseTree.Root, pos)
	if len(path) == 0 {
		return false
	}

	// Check if this variable is defined locally
	found := false
	eachDefinedVariableAtPos(path[len(path)-1], path[0].Range().From, func(name string) {
		if name == varName {
			found = true
		}
	})
	return found
}

// eachDefinedVariableAtPos calls f for each variable defined and visible at pos
// This is adapted from pkg/edit/complete/ns_helper.go:eachDefinedVariable
func eachDefinedVariableAtPos(n parse.Node, pos int, f func(string)) {
	if fn, ok := n.(*parse.Form); ok {
		eachDefinedVariableInForm(fn, f)
	}
	if pn, ok := n.(*parse.Primary); ok && pn.Type == parse.Lambda {
		for _, param := range pn.Elements {
			if varRef, ok := cmpd.StringLiteral(param); ok {
				_, name := eval.SplitSigil(varRef)
				f(name)
			}
		}
	}
	for _, ch := range parse.Children(n) {
		if ch.Range().From > pos {
			break
		}
		if pn, ok := ch.(*parse.Primary); ok && pn.Type == parse.Lambda {
			if pos >= pn.Range().To {
				continue
			}
		}
		eachDefinedVariableAtPos(ch, pos, f)
	}
}

// eachDefinedVariableInForm calls f for each variable defined in fn
// This is adapted from pkg/edit/complete/ns_helper.go:eachDefinedVariableInForm
func eachDefinedVariableInForm(fn *parse.Form, f func(string)) {
	if fn.Head == nil {
		return
	}
	switch head, _ := cmpd.StringLiteral(fn.Head); head {
	case "var":
		for _, arg := range fn.Args {
			if parse.SourceText(arg) == "=" {
				break
			}
			if varRef, ok := cmpd.StringLiteral(arg); ok {
				_, name := eval.SplitSigil(varRef)
				f(name)
			}
		}
	case "fn":
		if len(fn.Args) >= 1 {
			if name, ok := cmpd.StringLiteral(fn.Args[0]); ok {
				f(name + eval.FnSuffix)
			}
		}
	}
}

// getLocalVariableInfo returns information about a local variable if it exists
func getLocalVariableInfo(varName string) (string, error) {
	// For locally defined variables, we provide basic information
	// since they don't have full documentation like builtins
	return fmt.Sprintf("**%s** (local variable)\n\nLocally defined variable in current scope.", varName), nil
}

func newServer() *server {
	return &server{eval.NewEvaler(), make(map[lsp.DocumentURI]document)}
}

func handler(s *server) jsonrpc2.Handler {
	return routingHandler(map[string]method{
		"initialize":              s.initialize,
		"textDocument/didOpen":    convertMethod(s.didOpen),
		"textDocument/didChange":  convertMethod(s.didChange),
		"textDocument/hover":      convertMethod(s.hover),
		"textDocument/completion": convertMethod(s.completion),

		"textDocument/didClose": noop,
		// Required by spec.
		"initialized": noop,
		// Called by clients even when server doesn't advertise support:
		// https://microsoft.github.io/language-server-protocol/specification#workspace_didChangeWatchedFiles
		"workspace/didChangeWatchedFiles": noop,
	})
}

type method func(context.Context, json.RawMessage) (any, error)

func convertMethod[T any](f func(context.Context, T) (any, error)) method {
	return func(ctx context.Context, rawParams json.RawMessage) (any, error) {
		var params T
		if json.Unmarshal(rawParams, &params) != nil {
			return nil, errInvalidParams
		}
		return f(ctx, params)
	}
}

func noop(_ context.Context, _ json.RawMessage) (any, error) { return nil, nil }

type connKey struct{}

func routingHandler(methods map[string]method) jsonrpc2.Handler {
	return jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
		fn, ok := methods[req.Method]
		if !ok {
			return nil, errMethodNotFound
		}
		return fn(context.WithValue(ctx, connKey{}, conn), *req.Params)
	})
}

// Can be used within handler implementations to recover the connection stored
// in the Context.
func conn(ctx context.Context) *jsonrpc2.Conn { return ctx.Value(connKey{}).(*jsonrpc2.Conn) }

// Handler implementations. These are all called synchronously.

func (s *server) initialize(_ context.Context, _ json.RawMessage) (any, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    lsp.TDSyncKindFull,
			},
			CompletionProvider: &lsp.CompletionOptions{},
			HoverProvider:      &lsp.HoverOptions{},
		},
	}, nil
}

func (s *server) didOpen(ctx context.Context, params lsp.DidOpenTextDocumentParams) (any, error) {
	uri, content := params.TextDocument.URI, params.TextDocument.Text
	s.updateDocument(conn(ctx), uri, content)
	return nil, nil
}

func (s *server) didChange(ctx context.Context, params lsp.DidChangeTextDocumentParams) (any, error) {
	// ContentChanges includes full text since the server is only advertised to
	// support that; see the initialize method.
	uri, content := params.TextDocument.URI, params.ContentChanges[0].Text
	s.updateDocument(conn(ctx), uri, content)
	return nil, nil
}

func (s *server) hover(_ context.Context, params lsp.TextDocumentPositionParams) (any, error) {
	document, ok := s.documents[params.TextDocument.URI]
	if !ok {
		return nil, unknownDocument(params.TextDocument.URI)
	}
	pos := lspPositionToIdx(document.code, params.Position)

	p := np.Find(document.parseTree.Root, pos)
	// Try variable doc
	var primary *parse.Primary
	if p.Match(np.Store(&primary)) && primary.Type == parse.Variable {
		// Check if variable is locally defined (shadowing consideration)
		if isVariableShadowed(document.parseTree, pos, primary.Value) {
			// Variable is locally defined, provide local variable information
			markdown, err := getLocalVariableInfo(primary.Value)
			if err == nil {
				return lsp.Hover{Contents: lsp.MarkupContent{Kind: lsp.MKMarkdown, Value: markdown}}, nil
			}
		} else {
			// Variable is not locally defined, try global/builtin documentation
			markdown, err := doc.Source("$" + primary.Value)
			if err == nil {
				return lsp.Hover{Contents: lsp.MarkupContent{Kind: lsp.MKMarkdown, Value: markdown}}, nil
			}
		}
	}
	// Try command doc
	var expr np.SimpleExprData
	var form *parse.Form
	if p.Match(np.SimpleExpr(&expr, nil), np.Store(&form)) && form.Head == expr.Compound {
		// Check if command is locally defined (shadowing consideration)
		if isVariableShadowed(document.parseTree, pos, expr.Value) {
			// Command is locally defined, provide local function information
			markdown, err := getLocalVariableInfo(expr.Value)
			if err == nil {
				return lsp.Hover{Contents: lsp.MarkupContent{Kind: lsp.MKMarkdown, Value: markdown}}, nil
			}
		} else {
			// Command is not locally defined, try global/builtin documentation
			markdown, err := doc.Source(expr.Value)
			if err == nil {
				return lsp.Hover{Contents: lsp.MarkupContent{Kind: lsp.MKMarkdown, Value: markdown}}, nil
			}
		}
	}
	return nil, nil
}

func (s *server) completion(_ context.Context, params lsp.CompletionParams) (any, error) {
	document, ok := s.documents[params.TextDocument.URI]
	if !ok {
		return nil, unknownDocument(params.TextDocument.URI)
	}
	code := document.code
	pos := lspPositionToIdx(code, params.Position)
	result, err := complete.Complete(
		complete.CodeBuffer{
			Content: code,
			Dot:     pos,
		},
		s.evaler,
		complete.Config{},
	)
	if err != nil {
		return []lsp.CompletionItem{}, nil
	}

	lspItems := make([]lsp.CompletionItem, len(result.Items))
	lspRange := lspRangeFromRange(code, result.Replace)

	// Enhanced completion item kind determination with shadowing consideration
	getCompletionKind := func(name, contextName string) lsp.CompletionItemKind {
		switch contextName {
		case "command":
			// Check if command is locally defined
			if isVariableShadowed(document.parseTree, pos, name) {
				return lsp.CIKFunction // Local function
			}
			return lsp.CIKFunction // Global/builtin function
		case "variable":
			// Always return CIKVariable for compatibility, but use details to distinguish
			return lsp.CIKVariable
		case "argument":
			return lsp.CIKValue // Command arguments
		case "index":
			return lsp.CIKProperty // Array/map indices
		case "redir":
			return lsp.CIKFile // File paths for redirections
		default:
			return lsp.CIKText // Fallback for unknown types
		}
	}

	for i, item := range result.Items {
		kind := getCompletionKind(item.ToInsert, result.Name)

		// Enhanced completion item with additional information
		completionItem := lsp.CompletionItem{
			Label: item.ToInsert,
			Kind:  kind,
			TextEdit: &lsp.TextEdit{
				Range:   lspRange,
				NewText: item.ToInsert,
			},
		}

		// Add detail information for better UX
		switch result.Name {
		case "variable":
			if isVariableShadowed(document.parseTree, pos, item.ToInsert) {
				completionItem.Detail = "local variable"
			} else {
				completionItem.Detail = "global/builtin variable"
			}
		case "command":
			if isVariableShadowed(document.parseTree, pos, item.ToInsert) {
				completionItem.Detail = "local function"
			} else {
				completionItem.Detail = "global/builtin command"
			}
		case "argument":
			completionItem.Detail = "command argument"
		case "index":
			completionItem.Detail = "index key"
		case "redir":
			completionItem.Detail = "file path"
		}

		lspItems[i] = completionItem
	}
	return lspItems, nil
}

func (s *server) updateDocument(conn *jsonrpc2.Conn, uri lsp.DocumentURI, code string) {
	tree, err := parse.Parse(parse.Source{Name: string(uri), Code: code}, parse.Config{})
	s.documents[uri] = document{code, tree, err}
	go func() {
		// Convert the parse error to lsp.Diagnostic objects and publish them.
		entries := parse.UnpackErrors(err)
		diags := make([]lsp.Diagnostic, len(entries))
		for i, err := range entries {
			diags[i] = lsp.Diagnostic{
				Range:    lspRangeFromRange(code, err),
				Severity: lsp.DSError,
				Source:   "parse",
				Message:  err.Message,
			}
		}
		conn.Notify(context.Background(), "textDocument/publishDiagnostics",
			lsp.PublishDiagnosticsParams{URI: uri, Diagnostics: diags})
	}()
}

func unknownDocument(uri lsp.DocumentURI) error {
	return &jsonrpc2.Error{
		Code:    jsonrpc2.CodeInvalidParams,
		Message: fmt.Sprintf("unknown document: %v", uri),
	}
}

func lspRangeFromRange(s string, r diag.Ranger) lsp.Range {
	rg := r.Range()
	return lsp.Range{
		Start: lspPositionFromIdx(s, rg.From),
		End:   lspPositionFromIdx(s, rg.To),
	}
}

func lspPositionToIdx(s string, pos lsp.Position) int {
	var idx int
	walkString(s, func(i int, p lsp.Position) bool {
		idx = i
		return p.Line < pos.Line || (p.Line == pos.Line && p.Character < pos.Character)
	})
	return idx
}

func lspPositionFromIdx(s string, idx int) lsp.Position {
	var pos lsp.Position
	walkString(s, func(i int, p lsp.Position) bool {
		pos = p
		return i < idx
	})
	return pos
}

// Generates (index, lspPosition) pairs in s, stopping if f returns false.
func walkString(s string, f func(i int, p lsp.Position) bool) {
	var p lsp.Position
	lastCR := false

	for i, r := range s {
		if !f(i, p) {
			return
		}
		switch {
		case r == '\r':
			p.Line++
			p.Character = 0
		case r == '\n':
			if lastCR {
				// Ignore \n if it's part of a \r\n sequence
			} else {
				p.Line++
				p.Character = 0
			}
		case r <= 0xFFFF:
			// Encoded in UTF-16 with one unit
			p.Character++
		default:
			// Encoded in UTF-16 with two units
			p.Character += 2
		}
		lastCR = r == '\r'
	}
	f(len(s), p)
}
