package controllers

import "Projet-Forum/internal/middlewares"

var IndexHandlerGetBundle = middlewares.Join(indexHandlerGet, middlewares.Log, middlewares.UserCheck)
var IndexHandlerPutBundle = middlewares.Join(indexHandlerPut, middlewares.Log, middlewares.Guard)
var LoginHandlerGetBundle = middlewares.Join(loginHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var LoginHandlerPostBundle = middlewares.Join(loginHandlerPost, middlewares.Log, middlewares.OnlyVisitors)
var RegisterHandlerGetBundle = middlewares.Join(registerHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var RegisterHandlerPostBundle = middlewares.Join(registerHandlerPost, middlewares.Log, middlewares.OnlyVisitors)
var HomeHandlerGetBundle = middlewares.Join(homeHandlerGet, middlewares.Log, middlewares.Guard)
var LogHandlerGetBundle = middlewares.Join(logHandlerGet, middlewares.Log, middlewares.UserCheck)
var ConfirmHandlerGetBundle = middlewares.Join(confirmHandlerGet, middlewares.Log, middlewares.OnlyVisitors)
var LogoutHandlerGetBundle = middlewares.Join(logoutHandlerGet, middlewares.Log, middlewares.Guard)
var ErrorHandlerBundle = middlewares.Join(errorHandler, middlewares.Log, middlewares.UserCheck)

//var CreateCategoryPostBundle = middlewares.Join(createCategoryPost, middlewares.Log, middlewares.Guard)
//var CreateThreadPostBundle = middlewares.Join(createThreadPost, middlewares.Log, middlewares.Guard)
//var CreatePostPostBundle = middlewares.Join(createPostPost, middlewares.Log, middlewares.Guard)
//var CreateTagPostBundle = middlewares.Join(createTagPost, middlewares.Log, middlewares.Guard)
//var ThreadGetBundle = middlewares.Join(threadGet, middlewares.Log, middlewares.Guard)
//var TagGetBundle = middlewares.Join(tagGet, middlewares.Log, middlewares.Guard)
//var CategoryGetBundle = middlewares.Join(categoryGet, middlewares.Log, middlewares.Guard)
//var ProfileGetBundle = middlewares.Join(profileGet, middlewares.Log, middlewares.Guard)
