功能一：实时自动生成数据发送http请求api
1、init.conf中配置url地址及测试主体修改为true；配置渠道:场景标识
2、点击run.bat执行报文发送请求；所生成的数据均为自动产生
3、命中规则在log.txt中增量记录（可能存在重复命中）；使用Dr.bat对log.txt命中的规则进行去重，获取Dr_for_hitRules.txt规则命中结果
4、init.conf可以自动调用事件上报接口具体参加init.conf中参数注释
5、可打印每次请求的标记对应关系；开启是否打印请求报文及返回报文日志

功能二：Csv数据读取发送http请求api
1、配置如下参数
#读取csv数据发送http请求api
ifCsv = true
CsvName = ipaCsvData.csv
#是否一次性将文件读取，true为是
ifLoadCsvAll = true
#定义数据字段分割符(rune=int32，仅支持单字符（如,|&$%），除\r \n等占位符)
separator=$
#是否补充固定值字段
ifFieldsFixed = true
#因为时间日期字段比较特殊，因此需要配置时间日期字段所在的列是第几列，从0列开始，支持多列配置; 英文","隔开
isDate=0,15

2、按数据字段映射dataField.conf,列标从0开始映射，一一对应即可。
3、固定值传参fieldFixed.conf配置参数，详见样例

注：以上服务仅支持json报文，兼容新老平台（time参数和个别字段差异）