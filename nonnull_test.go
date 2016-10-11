package graphql_test

import (
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/flynn/graphql"
	"github.com/flynn/graphql/gqlerrors"
	"github.com/flynn/graphql/language/location"
	"github.com/flynn/graphql/testutil"
)

var syncError = "sync"
var nonNullSyncError = "nonNullSync"
var promiseError = "promise"
var nonNullPromiseError = "nonNullPromise"

var erroringData = map[string]interface{}{
	"sync": func() (interface{}, error) {
		return nil, errors.New(syncError)
	},
	"nonNullSync": func() (interface{}, error) {
		return nil, errors.New(nonNullSyncError)
	},
	"promise": func() (interface{}, error) {
		return nil, errors.New(promiseError)
	},
	"nonNullPromise": func() (interface{}, error) {
		return nil, errors.New(nonNullPromiseError)
	},
}

var nullingData = map[string]interface{}{
	"sync": func() interface{} {
		return nil
	},
	"nonNullSync": func() interface{} {
		return nil
	},
	"promise": func() interface{} {
		return nil
	},
	"nonNullPromise": func() interface{} {
		return nil
	},
}

var dataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DataType",
	Fields: graphql.Fields{
		"sync": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullSync": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"promise": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullPromise": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var nonNullTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: dataType,
})

func init() {
	erroringData["nest"] = func() interface{} {
		return erroringData
	}
	erroringData["nonNullNest"] = func() interface{} {
		return erroringData
	}
	erroringData["promiseNest"] = func() interface{} {
		return erroringData
	}
	erroringData["nonNullPromiseNest"] = func() interface{} {
		return erroringData
	}

	nullingData["nest"] = func() interface{} {
		return nullingData
	}
	nullingData["nonNullNest"] = func() interface{} {
		return nullingData
	}
	nullingData["promiseNest"] = func() interface{} {
		return nullingData
	}
	nullingData["nonNullPromiseNest"] = func() interface{} {
		return nullingData
	}

	dataType.AddFieldConfig("nest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
	dataType.AddFieldConfig("promiseNest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullPromiseNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
}

// nulls a nullable field that panics
func TestNonNull_NullsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: syncError,
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAComplexTreeOfNullableFieldsThatThrow(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 13},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldThrowsInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}

}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNull(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sync": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"promise": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAComplexTreeOfNullableFieldsThatReturnNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"nest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]interface{}{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]interface{}{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldReturnsNullInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected.Data, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !reflect.DeepEqual(expected.Errors, result.Errors) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}

func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldThrows(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldErrors(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   erroringData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldReturnsNull(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldResolvesNull(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if len(result.Errors) != len(expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
