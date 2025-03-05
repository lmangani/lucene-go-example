package main

import (
	"context"
	"fmt"
	"os"
	
	"github.com/geange/lucene-go/codecs/simpletext"
	"github.com/geange/lucene-go/core/document"
	"github.com/geange/lucene-go/core/index"
	"github.com/geange/lucene-go/core/search"
	"github.com/geange/lucene-go/core/store"
	index2 "github.com/geange/lucene-go/core/interface/index"
)

func main() {
	err := os.RemoveAll("data")
	if err != nil {
		panic(err)
	}

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

	{
		doc := document.NewDocument()
		doc.Add(document.NewTextField("a", "74", true))
		doc.Add(document.NewTextField("b", "86", true))
		doc.Add(document.NewTextField("c", "1237", true))
		docID, err := writer.AddDocument(ctx, doc)
		if err != nil {
			panic(err)
		}
		fmt.Println("add new document:", docID)
	}

	{
		doc := document.NewDocument()
		doc.Add(document.NewTextField("a", "74", true))
		doc.Add(document.NewTextField("b", "123", true))
		doc.Add(document.NewTextField("c", "789", true))

		docID, err := writer.AddDocument(ctx, doc)
		if err != nil {
			panic(err)
		}
		fmt.Println("add new document:", docID)
	}

	{
		doc := document.NewDocument()
		doc.Add(document.NewTextField("a", "741", true))
		doc.Add(document.NewTextField("b", "861", true))
		doc.Add(document.NewTextField("c", "12137", true))
		docID, err := writer.AddDocument(ctx, doc)
		if err != nil {
			panic(err)
		}
		fmt.Println("add new document:", docID)
	}
	
	// Force commit to ensure documents are indexed
	err = writer.Commit(ctx)
	if err != nil {
		panic(err)
	}
	
	// Run search query
	searchQuery()
}

func searchQuery() {
	dir, err := store.NewNIOFSDirectory("data")
	if err != nil {
		panic(err)
	}

	reader, err := index.OpenDirectoryReader(context.Background(), dir, nil, nil)
	if err != nil {
		panic(err)
	}
	// Don't close the reader here as it may cause issues in this version

	searcher, err := search.NewIndexSearcher(reader)
	if err != nil {
		panic(err)
	}

	query := search.NewTermQuery(index.NewTerm("a", []byte("74")))
	builder := search.NewBooleanQueryBuilder()
	builder.AddQuery(query, index2.OccurMust)
	booleanQuery, err := builder.Build()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	topDocs, err := searcher.SearchTopN(ctx, booleanQuery, 2)
	if err != nil {
		panic(err)
	}

	result := topDocs.GetScoreDocs()
	for _, scoreDoc := range result {
		docID := scoreDoc.GetDoc()
		document, err := reader.Document(ctx, docID)
		if err != nil {
			panic(err)
		}

		// Since we can't directly iterate document fields in version 0.0.2,
		// we'll retrieve the field values we know exist based on our document structure
		values := []string{
			fmt.Sprintf("Field a: %s", getFieldValue(document, "a")),
			fmt.Sprintf("Field b: %s", getFieldValue(document, "b")),
			fmt.Sprintf("Field c: %s", getFieldValue(document, "c")),
		}
		
		fmt.Printf("docId: %d, values: %v\n", scoreDoc.GetDoc(), values)
	}
}

// Helper function to get a field value by name
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
