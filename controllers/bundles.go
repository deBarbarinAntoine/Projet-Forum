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
