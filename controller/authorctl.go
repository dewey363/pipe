// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (C) 2017, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"net/http"

	"strings"

	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
	"github.com/vinta/pangu"
)

func showAuthorsAction(c *gin.Context) {
	blogAdmin := getBlogAdmin(c)
	dataModel := getDataModel(c)

	themeAuthors := []*ThemeAuthor{}
	authorModels := service.User.GetBlogUsers(blogAdmin.BlogID)
	for _, authorModel := range authorModels {
		author := &ThemeAuthor{
			Name:         authorModel.Name,
			URL:          getBlogURL(c) + util.PathAuthors + "/" + authorModel.Name,
			ArticleCount: authorModel.ArticleCount,
			AvatarURL:    authorModel.AvatarURLWithSize(210),
		}
		themeAuthors = append(themeAuthors, author)
	}

	dataModel["Authors"] = themeAuthors
	c.HTML(http.StatusOK, getTheme(c)+"/authors.html", dataModel)
}

func showAuthorArticlesAction(c *gin.Context) {
	authorName := strings.SplitAfter(c.Request.URL.Path, util.PathAuthors+"/")[1]
	author := service.User.GetUserByName(authorName)
	if nil == author {
		c.Status(404)

		return
	}

	page := c.GetInt("p")
	if 1 > page {
		page = 1
	}
	blogAdmin := getBlogAdmin(c)
	dataModel := getDataModel(c)
	session := util.GetSession(c)
	articleModels, pagination := service.Article.GetAuthorArticles(author.ID, page, blogAdmin.BlogID)
	articles := []*ThemeArticle{}
	for _, articleModel := range articleModels {
		themeTags := []*ThemeTag{}
		tagStrs := strings.Split(articleModel.Tags, ",")
		for _, tagStr := range tagStrs {
			themeTag := &ThemeTag{
				Title: tagStr,
				URL:   getBlogURL(c) + util.PathTags + "/" + tagStr,
			}
			themeTags = append(themeTags, themeTag)
		}

		authorModel := service.User.GetUser(articleModel.AuthorID)
		author := &ThemeAuthor{
			Name:      authorModel.Name,
			URL:       "http://localhost:5879/blogs/pipe/vanessa",
			AvatarURL: "https://img.hacpai.com/20170818zhixiaoyun.jpeg",
		}

		article := &ThemeArticle{
			ID:           articleModel.ID,
			Author:       author,
			CreatedAt:    articleModel.CreatedAt.Format("2006-01-02"),
			Title:        pangu.SpacingText(articleModel.Title),
			Tags:         themeTags,
			URL:          getBlogURL(c) + articleModel.Path,
			Topped:       articleModel.Topped,
			ViewCount:    articleModel.ViewCount,
			CommentCount: articleModel.CommentCount,
			ThumbnailURL: "https://img.hacpai.com/20170818zhixiaoyun.jpeg",
			Editable:     session.UID == authorModel.ID,
		}

		articles = append(articles, article)
	}
	dataModel["Articles"] = articles
	dataModel["Pagination"] = pagination
	dataModel["Author"] = author

	c.HTML(http.StatusOK, getTheme(c)+"/author-articles.html", dataModel)
}
