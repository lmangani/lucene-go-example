# Go Lucido
A simple demonstration of using [lucene-go](https://github.com/geange/lucene-go) for document indexing and search in Go.

## Overview

This project shows how to use lucene-go to create and search a simple document index. Lucene is a high-performance, full-featured text search engine library originally written in Java, and lucene-go provides similar functionality in Go.

## Installation

1. Ensure you have Go 1.22+ installed
2. Clone this repository
3. Install dependencies:

```bash
go mod tidy
```

## Usage

Run the example application:

```bash
go run main.go
```

## Code Examples

### Creating an Index and Adding Documents

```go
// Initialize directory and writer
dir, err := store.NewNIOFSDirectory("data")
if err != nil {
    panic(err)
}

codec := simpletext.NewCodec()
similarity, err := search.NewBM25Similarity()
if err != nil {
    panic(err)
}

config := index.NewIndexWriterConfig(codec, similarity)
ctx := context.Background()
writer, err := index.NewIndexWriter(ctx, dir, config)
if err != nil {
    panic(err)
}
defer func() {
    err := writer.Commit(ctx)
    if err != nil {
        panic(err)
    }
    writer.Close()
}()

// Add documents
doc := document.NewDocument()
doc.Add(document.NewTextField("title", "Example Document", true))
doc.Add(document.NewTextField("content", "This is the content of the document", true))
doc.Add(document.NewTextField("author", "John Doe", true))

docID, err := writer.AddDocument(ctx, doc)
if err != nil {
    panic(err)
}
fmt.Println("Added document with ID:", docID)

// Force commit to ensure documents are indexed
err = writer.Commit(ctx)
if err != nil {
    panic(err)
}
```

### Searching Documents

```go
// Open reader
reader, err := index.OpenDirectoryReader(context.Background(), dir, nil, nil)
if err != nil {
    panic(err)
}

// Create searcher
searcher, err := search.NewIndexSearcher(reader)
if err != nil {
    panic(err)
}

// Create a term query
query := search.NewTermQuery(index.NewTerm("title", []byte("Example")))

// Execute search
ctx := context.Background()
topDocs, err := searcher.SearchTopN(ctx, query, 10)
if err != nil {
    panic(err)
}

// Process results
result := topDocs.GetScoreDocs()
for _, scoreDoc := range result {
    docID := scoreDoc.GetDoc()
    document, err := reader.Document(ctx, docID)
    if err != nil {
        panic(err)
    }
    
    // Access fields
    title := getFieldValue(document, "title")
    content := getFieldValue(document, "content")
    author := getFieldValue(document, "author")
    
    fmt.Printf("Doc ID: %d\nTitle: %s\nContent: %s\nAuthor: %s\n\n", 
        docID, title, content, author)
}
```

### Helper Functions

```go
// Get a field value by name
func getFieldValue(doc *document.Document, fieldName string) string {
    field, found := doc.GetField(fieldName)
    if !found {
        return "not found"
    }
    value := field.Get()
    if strValue, ok := value.(string); ok {
        return strValue
    }
    return fmt.Sprintf("%v", value)
}
```

## Advanced Usage

### Boolean Queries

```go
// Create a boolean query with multiple conditions
builder := search.NewBooleanQueryBuilder()
builder.AddQuery(search.NewTermQuery(index.NewTerm("author", []byte("John"))), index2.OccurMust)
builder.AddQuery(search.NewTermQuery(index.NewTerm("title", []byte("Example"))), index2.OccurShould)
booleanQuery, err := builder.Build()
if err != nil {
    panic(err)
}

// Execute search with the boolean query
topDocs, err := searcher.SearchTopN(ctx, booleanQuery, 10)
```

## Notes

- This implementation uses lucene-go v0.0.2
- The SimpleText codec is used for human-readable index files (not recommended for production)
- Remember to explicitly commit changes to ensure documents are indexed

## Resources

- [lucene-go GitHub Repository](https://github.com/geange/lucene-go)
- [Apache Lucene Documentation](https://lucene.apache.org/core/9_0_0/index.html) (Java version, but concepts apply)
