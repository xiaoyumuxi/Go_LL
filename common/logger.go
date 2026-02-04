package common

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化 Logger
func InitLogger() {
	// 1. 配置 Lumberjack 日志切割
	writeSyncer := getLogWriter()

	// 2. 配置日志编码器
	encoder := getEncoder()

	// 3. 创建 Core
	// 同时输出到控制台和文件
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)

	// 4. 创建 Logger
	// AddCaller: 添加调用者信息 (文件名和行号)
	Logger = zap.New(core, zap.AddCaller())

	// 替换全局的 zap logger，方便后续直接使用 zap.L()
	zap.ReplaceGlobals(Logger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间格式为 ISO8601 (2006-01-02T15:04:05.000Z0700)
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 日志级别大写 (INFO, ERROR)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 使用 JSON 格式
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    10,               // 单个文件最大尺寸 (MB)
		MaxBackups: 5,                // 最多保留备份个数
		MaxAge:     30,               // 最多保留天数
		Compress:   false,            // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}
