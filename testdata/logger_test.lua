-- 引入 logger
logger = require("logger")

logger.logInfoLn("debug info")
log = logger.create("app.log")
log.setLevel("debug")
log.logInfoLn("创建日志")
log.logDebugLn("debug日志")

logger.logInfoLn("log level:",log.getLevel())