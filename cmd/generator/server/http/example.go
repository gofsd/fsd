package http

import "net/http"

const (
	emptyArray = "[]"
	emptyObject = `{}`
	firstUser = `{"id":1, "firstName": "Dima", "lastName": "Pohiba", "email": "madi.nickName@gmail.com", "login": "admin", "publicKey":""}`
	postFirstUser = `{"firstName": "Dima", "lastName": "Pohiba", "email": "madi.nickName@gmail.com", "login": "admin", "publicKey":""}`
	secondUser = `{"id": 2, "firstName": "Dima", "lastName": "Pohiba", "email": "madi.nickName@gmail.com", "login": "admin", "publicKey":""}`
	thirdUser = `{"id": 3, "firstName": "Dima", "lastName": "Pohiba", "email": "madi.nickName@gmail.com", "login": "admin", "publicKey":""}`
    postFirstArticle = `{"title": "first article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum", "userId": 1}`
    firstArticle = `{"id":1, "title": "first article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",  "userId": 1}`
    secondArticle = `{"id":2, "title": "second article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",  "userId": 1}`
    postSecondArticle = `{ "title": "second article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum", "userId": 1}`
    thirdArticle = `{"id": 3, "title": "three article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",  "userId": 1}`
    postThirdArticle = `{"title": "three article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum", "userId": 1}`
    patchThirdArticle = `{"title": "three article patched", "someNewfield":"some new value"}`
    patchedThirdArticle = `{"id": 3, "userId": 1, "title": "three article patched", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum", "someNewfield":"some new value"}`
    putThirdArticle = `{
					"userId": 1,
					"title": "three article",
					"body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
					  "info": {
						"posted": "now"
					  }
					}`
    newThirdArticle = `{"id": 3, "userId": 1, "title": "three article", "body": "laboriosam mollitia et enim quasi adipisci quia provident illum", 					  "info": {
						"posted": "now"
					  }}`
    findedArticleByTitle = `[{
					"id": 3,
					"userId": 1,
					"title": "three article",
					"body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
					  "info": {
						"posted": "now"
					  }
					}]`
    findedArticlesByTitleAndNestedField =`[{
					"userId": 1,
					"id": 3,
					"title": "three article",
					"body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
					  "info": {
						"posted": "now"
					  }
					}]`
    findedArticlesByIds = `[
          {
            "title": "first article",
            "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
            "id": 1,
			"userId": 1
          },
          {
            "title": "three article",
            "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
			"userId": 1,
            "info": {
              "posted": "now"
            },
            "id": 3
          }
        ]`
    paginateArticlesWithPageOneAndLimitOne = `[{
            "title": "first article",
            "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
            "id": 1,
			"userId": 1
          }]`
    linksWithArticles = `<http://127.0.0.1:3000/articles?_page=1&_limit=1>; rel="first", <http://127.0.0.1:3000/articles?_page=1&_limit=1>; rel="prev", <http://127.0.0.1:3000/articles?_page=3&_limit=1>; rel="next", <http://127.0.0.1:3000/articles?_page=3&_limit=1>; rel="last"`
    articlesTotalCount = `3`
    articlesPageTwo = "["+ secondArticle +"]"
    articlesForUserOne = `[
		  {
			"title": "first article",
			"body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
			"userId": 1,
			"id": 1
		  },
		  {
			"title": "second article",
			"body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
			"userId": 1,
			"id": 2
		  },
		  {
			"userId": 1,
			"title": "three article",
			"body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
			"info": {
			  "posted": "now"
			},
			"id": 3
		  }
	]`
    articlesForUserOneSorted = `[
  {
    "userId": 1,
    "title": "three article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "info": {
      "posted": "now"
    },
    "id": 3
  },
  {
    "title": "second article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "userId": 1,
    "id": 2
  },
  {
    "title": "first article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "userId": 1,
    "id": 1
  }
]`
    sliceArticlesForUserOne = `[
  {
    "userId": 1,
    "title": "three article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "info": {
      "posted": "now"
    },
    "id": 3
  }
]`
	articleWithAllOperators = `[
  {
    "userId": 1,
    "title": "three article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "info": {
      "posted": "now"
    },
    "id": 3
  }
]
`
	nestedArticlesInUserResp = `{
  "firstName": "Dima",
  "lastName": "Pohiba",
  "email": "madi.nickName@gmail.com",
  "login": "admin",
  "publicKey": "",
  "id": 1,
  "articles": [
    {
      "title": "first article",
      "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
      "userId": 1,
      "id": 1
    },
    {
      "title": "second article",
      "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
      "userId": 1,
      "id": 2
    },
    {
      "userId": 1,
      "title": "three article",
      "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
      "info": {
        "posted": "now"
      },
      "id": 3
    }
  ]
}`
	expandedArticlesInUserResp = `[
  {
    "title": "first article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "userId": 1,
    "id": 1,
    "user": {
      "firstName": "Dima",
      "lastName": "Pohiba",
      "email": "madi.nickName@gmail.com",
      "login": "admin",
      "publicKey": "",
      "id": 1
    }
  },
  {
    "title": "second article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "userId": 1,
    "id": 2,
    "user": {
      "firstName": "Dima",
      "lastName": "Pohiba",
      "email": "madi.nickName@gmail.com",
      "login": "admin",
      "publicKey": "",
      "id": 1
    }
  },
  {
    "userId": 1,
    "title": "three article",
    "body": "laboriosam mollitia et enim quasi adipisci quia provident illum",
    "info": {
      "posted": "now"
    },
    "id": 3,
    "user": {
      "firstName": "Dima",
      "lastName": "Pohiba",
      "email": "madi.nickName@gmail.com",
      "login": "admin",
      "publicKey": "",
      "id": 1
    }
  }
]`
)

//TestData some test data
var TestData = []EndpointTestData{
	{
		Name:     "Get an empty list of users",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "users",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get an empty list of articles",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "articles",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get an empty list of comments",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "comments",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get an empty list of pictures",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "pictures",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get an empty list of books",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "books",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get empty video list",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "video",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get empty audio list",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "audio",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get an empty list of articles 7",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyArray),
		},
		Request: Request{
			URL:        "articles",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Post user 1",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(firstUser),
		},
		Request: Request{
			URL:        "users",
			MethodType: http.MethodPost,
			Body: []byte(postFirstUser),
		},
	},
	{
		Name:     "Post user 2",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(secondUser),
		},
		Request: Request{
			URL:        "users",
			MethodType: http.MethodPost,
			Body: []byte(postFirstUser),
		},
	},
	{
		Name:     "Post user 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(thirdUser),
		},
		Request: Request{
			URL:        "users",
			MethodType: http.MethodPost,
			Body: []byte(postFirstUser),
		},
	},
	{
		Name:     "Post articles 1",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(firstArticle),
		},
		Request: Request{
			URL:        "articles",
			MethodType: http.MethodPost,
			Body: []byte(postFirstArticle),
		},
	},
	{
		Name:     "Post articles 2",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(secondArticle),
		},

		Request: Request{
			URL:        "articles",
			MethodType: http.MethodPost,
			Body: []byte(postSecondArticle),
		},
	},
	{
		Name:     "Post articles 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(thirdArticle),
		},

		Request: Request{
			URL:        "articles",
			MethodType: http.MethodPost,
			Body: []byte(postThirdArticle),
		},
	},
	{
		Name:     "Delete articles 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(emptyObject),
		},
		Request: Request{
			URL:        "articles/3",
			MethodType: http.MethodDelete,
		},
	},
	{
		Name:     "Post again articles 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(thirdArticle),
		},

		Request: Request{
			URL:        "articles",
			MethodType: http.MethodPost,
			Body: []byte(postThirdArticle),
		},
	},
	{
		Name:     "Patch articles 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(patchedThirdArticle),
		},

		Request: Request{
			URL:        "articles/3",
			MethodType: http.MethodPatch,
			Body:[]byte(patchThirdArticle),
		},
	},
	{
		Name:     "Put articles 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(newThirdArticle),
		},
		Request: Request{
			URL:        "articles/3",
			MethodType: http.MethodPut,
			Body: []byte(putThirdArticle),
		},
	},
	{
		Name:     "Search articles with title 'three article'",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(findedArticleByTitle),
		},
		Request: Request{
			URL:        "articles?title=three%20article",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Search articles with title 'three article' and nested structure info.posted",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(findedArticlesByTitleAndNestedField),
		},
		Request: Request{
			URL:        "articles?title=three%20article&info.posted=now",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Search articles with id 1 and 3",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(findedArticlesByIds),
		},
		Request: Request{
			URL:        "articles?id=1&id=3",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get articles page one with limit one",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(paginateArticlesWithPageOneAndLimitOne),
		},
		Request: Request{
			URL:        "articles?_page=1&_limit=1",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get articles page two with limit one",
		ExpectedResponse: ExpectedResponse{
			Headers: map[string][]string{
				"Link": []string{linksWithArticles},
				"X-Total-Count": []string{articlesTotalCount},

			},
			Body: []byte(articlesPageTwo),
		},
		Request: Request{
			URL:        "articles?_page=2&_limit=1",
			MethodType: http.MethodGet,
		},
	},

	{
		Name:     "Get articles for user 1",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(articlesForUserOne),
		},
		Request: Request{
			URL:        "users/1/articles",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get articles for user 1 sorted by info.posted and id asc, desc",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(articlesForUserOneSorted),
		},
		Request: Request{
			URL:        "users/1/articles?_sort=info.posted,id&_order=asc,desc",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Example all operator usage on articles for user one",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(articleWithAllOperators),
		},
		Request: Request{
			URL:        "users/1/articles?title_like=three&id_gte=1&id_lte=3&id_ne=1&q=adip",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get nested response",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(nestedArticlesInUserResp),
		},
		Request: Request{
			URL:        "users/1?_embed=articles",
			MethodType: http.MethodGet,
		},
	},
	{
		Name:     "Get expanded response",
		ExpectedResponse: ExpectedResponse{
			Body: []byte(expandedArticlesInUserResp),
		},
		Request: Request{
			URL:        "articles?_expand=user",
			MethodType: http.MethodGet,
		},
	},
}