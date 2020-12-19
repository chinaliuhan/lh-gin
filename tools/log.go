package tools

import (
	"go.uber.org/zap"
	"log"
)

type LogUtil struct {
	loggerHandler *zap.Logger
	logLevel      int
}

func NewLogUtil() *LogUtil {

	//直接使用默认方式启动,不自定义配置
	//logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()

	//encoder := zap.NewProductionEncoderConfig()
	//encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	//////使用自定义配置启动
	//encoderConfig := zapcore.EncoderConfig{
	//	MessageKey:       "lh-gin", //输入信息的key名
	//	LevelKey:         "aaa", //输出日志级别的key名
	//	TimeKey:          "bbb", //输出时间的key名
	//	NameKey:          "ccc",
	//	CallerKey:        "",
	//	FunctionKey:      "",
	//	StacktraceKey:    "",
	//	LineEnding:       ",",  //每行的分隔符。基本zapcore.DefaultLineEnding 即"\n"
	//	EncodeLevel:      nil, //基本zapcore.LowercaseLevelEncoder。将日志级别字符串转化为小写
	//	EncodeTime:       nil, //输出的时间格式
	//	EncodeDuration:   nil, //一般zapcore.SecondsDurationEncoder,执行消耗的时间转化成浮点型的秒
	//	EncodeCaller:     nil, //一般zapcore.ShortCallerEncoder，以包/文件:行号 格式化调用堆栈
	//	EncodeName:       nil,
	//	ConsoleSeparator: "",
	//}
	//
	//config := zap.Config{
	//	Level:             zap.AtomicLevel{}, //日志级别，
	//	Development:       false,             //是否是开发环境。如果是开发模式，对DPanicLevel进行堆栈跟踪
	//	DisableCaller:     false,             //禁止使用调用函数的文件名和行号来注释日志。默认进行注释日志
	//	DisableStacktrace: false,             //是否禁用堆栈跟踪捕获。默认对Warn级别以上和生产error级别以上的进行堆栈跟踪。
	//	Sampling:          nil,
	//	Encoding:          "",                      //编码类型，目前两种json 和 console【按照空格隔开】,常用json
	//	EncoderConfig:     zapcore.EncoderConfig{}, //生成格式的一些配置
	//	EncoderConfig:     encoderConfig,           //生成格式的一些配置
	//	OutputPaths:       nil,                     //日志写入文件的地址,如果想控制台输出，OutputPaths和ErrorOutputPaths不能配置为文件地址，而应该改为stdout。
	//	ErrorOutputPaths:  nil,                     //将系统内的error记录到文件的地址,如果想控制台输出，OutputPaths和ErrorOutputPaths不能配置为文件地址，而应该改为stdout。
	//	InitialFields:     nil,                     //加入一些初始的字段数据，比如项目名
	//}

	//config := zap.Config{
	//	EncoderConfig: zapcore.EncoderConfig{
	//		EncodeTime: zapcore.ISO8601TimeEncoder,
	//	},
	//}
	//encoder := zap.NewProductionEncoderConfig()
	//encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//logger, err := config.Build()
	if err != nil {
		log.Println("日志系统启动失败,构造函数实例化失败")
		return &LogUtil{}
	}
	defer logger.Sync()
	return &LogUtil{loggerHandler: logger}
}

/**
"debug":  zapcore.DebugLevel,
"info":   zapcore.InfoLevel,
"warn":   zapcore.WarnLevel,
"error":  zapcore.ErrorLevel,
"dpanic": zapcore.DPanicLevel,
"panic":  zapcore.PanicLevel,
"fatal":  zapcore.FatalLevel,
*/
const (
	LOG_DEBUG_LEVEL = iota
	LOG_INFO_LEVEL
	LOG_WARN_LEVEL
	LOG_ERROR_LEVEL
	LOG_PANIC_LEVEL
	LOG_FATAL_LEVEL
)

func (r *LogUtil) SetLevel(level int) *LogUtil {
	r.logLevel = level
	return r
}

func (r *LogUtil) SugarPrint(params ...interface{}) {
	switch r.logLevel {
	case LOG_DEBUG_LEVEL:
		r.loggerHandler.Sugar().Debug(params...)
	case LOG_INFO_LEVEL:
		r.loggerHandler.Sugar().Info(params...)
	case LOG_WARN_LEVEL:
		r.loggerHandler.Sugar().Warn(params...)
	case LOG_ERROR_LEVEL:
		r.loggerHandler.Sugar().Error(params...)
	case LOG_PANIC_LEVEL:
		r.loggerHandler.Sugar().Panic(params...)
	case LOG_FATAL_LEVEL:
		r.loggerHandler.Sugar().Fatal(params...)
	default:
		r.loggerHandler.Sugar().Info(params...)
	}
}

//func (r *LogUtil) Debug(text string, field ...zap.Field) error {
func (r *LogUtil) Debug(text ...interface{}) {
	//r.loggerHandler.Debug(text, field...)
	r.SugarPrint(text...)
	r.SetLevel(LOG_DEBUG_LEVEL)
}

//func (r *LogUtil) Info(text string, field ...zap.Field) error {
func (r *LogUtil) Info(text ...interface{}) {
	//r.loggerHandler.Info(text, field...)
	r.SugarPrint(text...)
	r.SetLevel(LOG_INFO_LEVEL)
}

//func (r *LogUtil) Warning(text string, field ...zap.Field) error {
func (r *LogUtil) Warning(text ...interface{}) {
	//r.loggerHandler.Warn(text, field...)
	r.SugarPrint(text...)
	r.SetLevel(LOG_WARN_LEVEL)
}

//func (r *LogUtil) Error(text string, field ...zap.Field) error {
func (r *LogUtil) Error(text ...interface{}) {
	//r.loggerHandler.Error(text, field...)
	r.SugarPrint(text...)
	r.SetLevel(LOG_ERROR_LEVEL)
}

//func (r *LogUtil) Fatal(text string, field ...zap.Field) error {
func (r *LogUtil) Fatal(text ...interface{}) {
	//r.loggerHandler.Fatal(text, field...)
	r.SugarPrint(text...)
	r.SetLevel(LOG_FATAL_LEVEL)
}

//func (r *LogUtil) Panic(text string, field ...zap.Field) error {
func (r *LogUtil) Panic(text ...interface{}) {
	//r.loggerHandler.Panic(text, field...)
	r.SugarPrint(text...)
	r.SetLevel(LOG_PANIC_LEVEL)
}
