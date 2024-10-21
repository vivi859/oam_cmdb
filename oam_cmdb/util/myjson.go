package util

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

const (
	TAG_TIME_FORMAT = "format"
	TAG_TIME_LOCATE = "locale"
)

//var JsonCustom = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	defaultFormat           = "2006-01-02 15:04:05"
	defaultLocale           = time.Local
	fieldBeginWithLowercase = true
)

//定义几个宽泛的日期匹配正则
var supportedFormatMap = map[string]string{
	`\d{4}\-\d{2}\-\d{2}\s\d{2}:\d{2}:\d{2}`: defaultFormat,
	`\d{4}/\d{2}\-\d{2}\s\d{2}:\d{2}:\d{2}`:  "2006/01/02 15:04:05",
	`\d{14}`:                                 "20060102150405"}

//新增一个可解析的日期格式.符合正则pattern的值将使用time_format解析
func AddSupportedTimeFormat(pattern string, time_format string) {
	supportedFormatMap[pattern] = time_format
}

func init() {
	jsoniter.RegisterExtension(&MyJsonIterExtension{})

	//自定义时间解析
	jsoniter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		val := iter.ReadString()
		var t time.Time
		var err error
		var parsed = false

		if strings.Contains(val, "T") {
			t, err = time.Parse(time.RFC3339, val)
			if err != nil {
				iter.Error = err
				return
			}
		} else {
			//尝试使用已知格式解析
			for k, v := range supportedFormatMap {
				if ok, _ := regexp.MatchString(k, val); ok {
					t, err = time.ParseInLocation(v, val, defaultLocale)
					if err == nil {
						parsed = true
						break
					}
				}
			}
			if !parsed {
				iter.Error = errors.New("无法解析" + val + ",不支持的时间格式")
				return
			}
		}

		*((*time.Time)(ptr)) = t
	})
}

type MyJsonIterExtension struct {
	jsoniter.DummyExtension
}

/*
定制json输出策略(坑爹的go真是麻烦).

1. 修改属性名为首字母小写,更符合json通用规则

2. 时间类型,默认格式修改为yyyy-MM-dd HH:mm:ss,同时可用tag format指定格式
*/
func (extension *MyJsonIterExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		//属性名统一改为json通用的首字母小写0
		if fieldBeginWithLowercase {
			tag, hastag := binding.Field.Tag().Lookup("json")
			if hastag {
				tagParts := strings.Split(tag, ",")
				if tagParts[0] == "-" {
					continue
				}
			}
			alias := ToFirstLetterLower(binding.Field.Name())
			binding.ToNames = []string{alias}
			//binding.FromNames = []string{binding.Field.Name(), alias, strings.ToLower(binding.Field.Name())}
		}

		//定制时间格式
		var typeErr error
		var isPtr bool
		typeName := binding.Field.Type().String()
		if typeName == "time.Time" {
			isPtr = false
		} else if typeName == "*time.Time" {
			isPtr = true
		} else {
			continue
		}

		var format string
		var isSpecifiedFormat = false
		formatTag := binding.Field.Tag().Get(TAG_TIME_FORMAT)
		if formatTag != "" {
			format = formatTag
			isSpecifiedFormat = true
		} else {
			format = defaultFormat
		}

		var locale *time.Location
		if localeTag := binding.Field.Tag().Get(TAG_TIME_LOCATE); localeTag != "" {
			loc, err := time.LoadLocation(localeTag)
			if err != nil {
				typeErr = err
			} else {
				locale = loc
			}
		} else {
			locale = defaultLocale
		}

		binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			if typeErr != nil {
				stream.Error = typeErr
				return
			}

			var tp *time.Time
			if isPtr {
				tpp := (**time.Time)(ptr)
				tp = *(tpp)
			} else {
				tp = (*time.Time)(ptr)
			}

			if tp != nil {
				lt := tp.In(locale)
				str := lt.Format(format)
				stream.WriteString(str)
			} else {
				stream.Write([]byte("null"))
			}
		}}
		//如果tag里指定的时间格式加特定解析,没有则用上面注册的宽格式解析
		if isSpecifiedFormat {
			binding.Decoder = &funcDecoder{fun: func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
				if typeErr != nil {
					iter.Error = typeErr
					return
				}

				str := iter.ReadString()
				var t *time.Time
				if str != "" {
					var err error
					tmp, err := time.ParseInLocation(format, str, locale)
					if err != nil {
						iter.Error = err
						return
					}
					t = &tmp
				} else {
					t = nil
				}

				if isPtr {
					tpp := (**time.Time)(ptr)
					*tpp = t
				} else {
					tp := (*time.Time)(ptr)
					if tp != nil && t != nil {
						*tp = *t
					}
				}
			}}
		}

	}
}

type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}
