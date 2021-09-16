/*
Package server implements a DocumentServer ready to serve a document to any number of clients. Each client can edit
the document simultaneously and the server will perform all the necessary transformations, broadcasting them to all
clients.

A custom StateStore can be implemented to use a database or the included MemoryStateStore can be used to store
everything in memory.
*/
package server
