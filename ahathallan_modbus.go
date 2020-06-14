package main

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/boltdb/bolt"
    "github.com/goburrow/modbus"
    "github.com/labstack/echo"
    _ "github.com/lib/pq"
    "log"
    "math"
    "net/http"
    "os"
    "os/exec"
    "reflect"
    "strconv"
    "time"

    uuid "satori/go.uuid"
    //simplejson "github.com/bitly/go-simplejson"
)

//日志对象
func LogInfo() *log.Logger {
    //创建文件对象, 日志的格式为当前时间2006-01-02 15:04:05.log;据说是golang的诞生时间，固定写法
    timeString := time.Now().Format("2006-01-02")
    file := "./A090_Atom_Modbus/log/"+timeString+".log"

    logFile,err := os.OpenFile(file,os.O_RDWR|os.O_CREATE|os.O_APPEND,0766)
    if(err != nil){
        panic(err)
    }

    //创建一个Logger, 参数1：日志写入的文件, 参数2：每条日志的前缀；参数3：日志属性
    return log.New(logFile,"logpre_",log.Lshortfile)
}
//日志对象
var mylooger *log.Logger
//基础间隔
var time_second=time.Duration(5)

var IOT_CONF Iot_conf

//配置文件对象
type Iot_conf struct {
    Mqtt_ip                     string
    Web_socket                  string
    Web_folder                  string
    Time_second                 int
    Mqtt_client_id              string
    Mqtt_client_watch_id        string
    Bolt_db_name                string
    Post_gres_userName          string
    Post_gres_password          string
    Post_gres_ipAddrees         string
    Post_gres_port              int
    Post_gres_dbName            string
    Base_name                   string
    Begin_time                  string
    Map_json                    string
}

//获取配置文件
func get_conf() Iot_conf{
    var iot_conf Iot_conf
    //var favorites songs
    var path string = "./A090_Atom_Modbus/conf_edge.toml"
    if _, err := toml.DecodeFile(path, &iot_conf); err != nil {
        mylooger.Output(2,err.Error())
    }

    return iot_conf
}

var MAIN_ Main_=Main_{}

type Main_ struct{

}

//boltdb操作-------------↓↓------------------------boltdb操作-------------↓↓------------------------boltdb操作-------------↓↓------------------------

type App_bolt_opt struct{
    db_file_name string
    bucket_name string
}

func new_App_bolt_opt(db_name string) *App_bolt_opt {
    var result=new(App_bolt_opt)

    result.db_file_name=IOT_CONF.Bolt_db_name

    result.bucket_name=db_name

    result.creat_db()

    return result
}

func new_App_bolt_opt_query(db_name string) *App_bolt_opt {
    var result=new(App_bolt_opt)

    result.db_file_name=IOT_CONF.Bolt_db_name

    result.bucket_name=db_name

    return result
}

func (core_bolt_opt *App_bolt_opt) creat_db(){

    //fmt.Printf("App_bolt_opt creat开始写入数据库%s:%s\n",core_bolt_opt.bucket_name,"================================")

    db, err:= bolt.Open(core_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
    }
    defer func(){
        //fmt.Printf("App_bolt_opt creat自动关闭数据库%s:%s\n",core_bolt_opt.bucket_name)
        db.Close()
    }()

    err = db.Update(func(tx *bolt.Tx) error {

        //判断要创建的表是否存在
        b := tx.Bucket([]byte(core_bolt_opt.bucket_name))
        if b == nil {

            //创建叫"MyBucket"的表
            _, err := tx.CreateBucket([]byte(core_bolt_opt.bucket_name))
            if err != nil {
                //也可以在这里对表做插入操作
                fmt.Printf("App_bolt_opt create错误1:%s:%s\n",core_bolt_opt.bucket_name,err)
                mylooger.Output(2,err.Error())
            }
        }

        //一定要返回nil
        return nil
    })

    //更新数据库失败
    if err != nil {
        fmt.Printf("App_bolt_opt create错误1:%s:%s\n",core_bolt_opt.bucket_name,err)
        mylooger.Output(2,err.Error())
    }

    if db!=nil {
        //fmt.Printf("App_bolt_opt create手动关闭数据库%s:%s\n",core_bolt_opt.bucket_name)
        db.Close()
    }
}

func (web_bolt_opt *App_bolt_opt) save_db_by_key(key string,web_info interface{}) {

    d_json, _ := json.Marshal(web_info)

    db, err:= bolt.Open(web_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
        fmt.Printf("Web_bolt_opt错误1:%s:%s\n",web_bolt_opt.bucket_name,err)
    }
    defer func(){
        // 获取异常信息
        if err:=recover();err!=nil{
            //  输出异常信息
            //fmt.Printf("Web_bolt_opt自动关闭数据库%s:%s\n", web_bolt_opt.bucket_name,err)
            errStr:= fmt.Sprint("watch_error:",err)
            mylooger.Output(2,errStr)
        }

        if db!=nil {
            //fmt.Printf("Web_bolt_opt自动关闭数据库%s:%s\n", web_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    err = db.Update(func(tx *bolt.Tx) error {

        //取出叫"MyBucket"的表
        b := tx.Bucket([]byte(web_bolt_opt.bucket_name))

        //往表里面存储数据
        if b != nil {
            //插入的键值对数据类型必须是字节数组
            err := b.Put([]byte(key), []byte(d_json))
            if err != nil {
                mylooger.Output(2,err.Error())
                fmt.Printf("Web_bolt_opt错误2:%s:%s\n",web_bolt_opt.bucket_name,err)
            }
        }

        //一定要返回nil
        return nil
    })

    //更新数据库失败
    if err != nil {
        mylooger.Output(2,err.Error())
        fmt.Printf("Web_bolt_opt错误3:%s:%s\n",web_bolt_opt.bucket_name,err)
    }

    if db!=nil {
        db.Close()
    }
}

func (web_bolt_opt *App_bolt_opt) save_db(key string,web_info interface{}) {

    d_json, _ := json.Marshal(web_info)

    fmt.Printf("Web_bolt_opt开始写入数据库%s:%s\n",web_bolt_opt.bucket_name,key)

    db, err:= bolt.Open(web_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
        fmt.Printf("Web_bolt_opt错误1:%s:%s\n",web_bolt_opt.bucket_name,err)
    }
    defer func(){
        // 获取异常信息
        if err:=recover();err!=nil{
            //  输出异常信息
            //fmt.Printf("Web_bolt_opt自动关闭数据库%s:%s\n", web_bolt_opt.bucket_name,err)
            errStr:= fmt.Sprint("watch_error:",err)
            mylooger.Output(2,errStr)
        }

        if db!=nil {
            //fmt.Printf("Web_bolt_opt自动关闭数据库%s:%s\n", web_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    err = db.Update(func(tx *bolt.Tx) error {

        //取出叫"MyBucket"的表
        b := tx.Bucket([]byte(web_bolt_opt.bucket_name))

        //往表里面存储数据
        if b != nil {
            //插入的键值对数据类型必须是字节数组
            err := b.Put([]byte(key), []byte(d_json))
            if err != nil {
                mylooger.Output(2,err.Error())
                fmt.Printf("Web_bolt_opt错误2:%s:%s\n",web_bolt_opt.bucket_name,err)
            }
        }

        //一定要返回nil
        return nil
    })

    //更新数据库失败
    if err != nil {
        mylooger.Output(2,err.Error())
        fmt.Printf("Web_bolt_opt错误3:%s:%s\n",web_bolt_opt.bucket_name,err)
    }

    fmt.Printf("Web_bolt_opt结束写入数据库%s:%s\n",web_bolt_opt.bucket_name,key)

    if db!=nil {
        fmt.Printf("Web_bolt_opt手动关闭数据库%s:%s\n",web_bolt_opt.bucket_name,key)
        db.Close()
    }
}

func (core_bolt_opt *App_bolt_opt) get_key_db (key string,web_info interface{}){
    db, err:= bolt.Open(core_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
    }
    defer func(){
        if db!=nil {
            //fmt.Printf("Web_bolt_opt get_key_db自动关闭数据库%s:%s\n", web_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    //4.查看表数据
    err = db.View(func(tx *bolt.Tx) error {

        //取出叫"MyBucket"的表
        b := tx.Bucket([]byte(core_bolt_opt.bucket_name))

        //往表里面存储数据
        if b != nil {

            data := b.Get([]byte(key))

            if data!=nil {
                _ = json.Unmarshal(data, &web_info)
            }
        }

        //一定要返回nil
        return nil
    })

    //查询数据库失败
    if err != nil {
        fmt.Printf("Web_bolt_opt get_key_db错误3:%s:%s\n",core_bolt_opt.bucket_name,err)
        mylooger.Output(2,err.Error())
    }
}

func (core_bolt_opt *App_bolt_opt) get_last_db(t reflect.Type) interface{}{

    db, err:= bolt.Open(core_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
    }
    defer func(){
        if db!=nil {
            //fmt.Printf("Core_bolt_opt_get_time_db自动关闭数据库%s:%s\n", core_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    var i_B_Device interface{}
    // Assume our events bucket exists and has RFC3339 encoded time keys.
    err = db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(core_bolt_opt.bucket_name)).Cursor()

        // Iterate over the 90's.
        _, v := c.Last()
        h := reflect.New(t).Interface()
        _ = json.Unmarshal(v, h)
        i_B_Device = h

        //一定要返回nil
        return nil
    })

    //查询数据库失败
    if err != nil {
        fmt.Printf("App_bolt_opt get_time_db错误1:%s:%s\n",core_bolt_opt.bucket_name,err)
        mylooger.Output(2,err.Error())
    }

    return i_B_Device
}

func (core_bolt_opt *App_bolt_opt) get_all_db(t reflect.Type) []interface{}{

    db, err:= bolt.Open(core_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
    }
    defer func(){
        if db!=nil {
            //fmt.Printf("Core_bolt_opt_get_time_db自动关闭数据库%s:%s\n", core_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    var i_d_lst []interface{}
    //4.查看表数据
    err = db.View(func(tx *bolt.Tx) error {

        //取出叫"MyBucket"的表
        b := tx.Bucket([]byte(core_bolt_opt.bucket_name))

        //往表里面存储数据
        if b != nil {
            b.ForEach(func(k, v []byte) error {
                if v != nil {
                    h := reflect.New(t).Interface()
                    _ = json.Unmarshal(v, &h)
                    var i_B_Device interface{}
                    i_B_Device = h

                    i_d_lst = append(i_d_lst, i_B_Device)
                }
                return nil
            })
        }

        //一定要返回nil
        return nil
    })

    //查询数据库失败
    if err != nil {
        fmt.Printf("Web_bolt_opt get_key_db错误3:%s:%s\n",core_bolt_opt.bucket_name,err)
        mylooger.Output(2,err.Error())
    }
    return i_d_lst
}

func (core_bolt_opt *App_bolt_opt) get_time_db(time_start time.Time,time_end time.Time,t reflect.Type) []interface{}{

    db, err:= bolt.Open(core_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
    }
    defer func(){
        if db!=nil {
            //fmt.Printf("Core_bolt_opt_get_time_db自动关闭数据库%s:%s\n", core_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    var i_d_lst []interface{}
    // Assume our events bucket exists and has RFC3339 encoded time keys.
    err = db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(core_bolt_opt.bucket_name)).Cursor()

        // Our time range spans the 90's decade.
        min := []byte(time_start.Format("2006-01-02 15:04:05"))
        max := []byte(time_end.Format("2006-01-02 15:04:05"))

        // Iterate over the 90's.
        for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
            h := reflect.New(t).Interface()
            _ = json.Unmarshal(v, h)
            var i_B_Device interface{}
            i_B_Device=h

            i_d_lst=append(i_d_lst, i_B_Device)
        }

        //一定要返回nil
        return nil
    })

    //查询数据库失败
    if err != nil {
        fmt.Printf("App_bolt_opt get_time_db错误1:%s:%s\n",core_bolt_opt.bucket_name,err)
        mylooger.Output(2,err.Error())
    }

    return i_d_lst
}

func (core_bolt_opt *App_bolt_opt) delete_db(key string) {

    //fmt.Printf("App_bolt_opt开始写入数据库%s:%s\n",core_bolt_opt.bucket_name)

    db, err:= bolt.Open(core_bolt_opt.db_file_name, 0600, nil)
    if err != nil {
        mylooger.Output(2,err.Error())
        fmt.Printf("App_bolt_opt错误1:%s:%s\n",core_bolt_opt.bucket_name,err)
    }
    defer func(){
        if db!=nil {
            //fmt.Printf("App_bolt_opt自动关闭数据库%s:%s\n", core_bolt_opt.bucket_name)
            db.Close()
        }
    }()

    err = db.Update(func(tx *bolt.Tx) error {

        //取出叫"MyBucket"的表
        b := tx.Bucket([]byte(core_bolt_opt.bucket_name))

        //往表里面存储数据
        if b != nil {
            //插入的键值对数据类型必须是字节数组
            err := b.Delete([]byte(key))
            if err != nil {
                mylooger.Output(2,err.Error())
                fmt.Printf("App_bolt_opt错误2:%s:%s\n",core_bolt_opt.bucket_name,err)
            }
        }

        //一定要返回nil
        return nil
    })

    //更新数据库失败
    if err != nil {
        mylooger.Output(2,err.Error())
        fmt.Printf("App_bolt_opt错误3:%s:%s\n",core_bolt_opt.bucket_name,err)
    }

    //fmt.Printf("App_bolt_opt结束写入数据库%s:%s\n",core_bolt_opt.bucket_name)

    if db!=nil {
        //fmt.Printf("App_bolt_opt手动关闭数据库%s:%s\n",core_bolt_opt.bucket_name)
        db.Close()
    }
}

//boltdb操作------------↑↑------------------------boltdb操作------------↑↑------------------------boltdb操作------------↑↑------------------------

// ★查询基站数据
func modbus_islink_restful(c echo.Context) error {
    return c.JSON(http.StatusOK, _MODBUS_MAIN.Is_link)
}

// ★查询基站数据
func modbus_link_restful(c echo.Context) error {
    _MODBUS_MAIN.link_modbus()
    return c.String(http.StatusOK, "")
}

// ★查询基站数据
func modbus_dis_link_restful(c echo.Context) error {
    _MODBUS_MAIN.dis_link_modbus()
    return c.String(http.StatusOK, "")
}

// ★保存Modbus方案
func modbus_scheme_save_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    modbus_scheme := Modbus_scheme{}

    json.Unmarshal(buf[0:n], &modbus_scheme)

    _MODBUS_MAIN.save_modbus_schedual(modbus_scheme)
    return c.String(http.StatusOK, "")
}

// ★获取Modbus方案列表
func modbus_scheme_get_lst_restful(c echo.Context) error {
    lst_modbus_shceme:= _MODBUS_MAIN.modbus_scheme_get_lst()

    return c.JSON(http.StatusOK, lst_modbus_shceme)
}

// ★获取Modbus方案列表
func modbus_scheme_get_restful(c echo.Context) error {
    modbus_scheme_id := c.Param("modbus_scheme_id")

    modbus_shceme := _MODBUS_MAIN.modbus_scheme_get(modbus_scheme_id)

    return c.JSON(http.StatusOK, modbus_shceme)
}

// ★删除Modbus方案
func modbus_scheme_delete_restful(c echo.Context) error {
    modbus_scheme_id := c.Param("modbus_scheme_id")

    _MODBUS_MAIN.modbus_scheme_delete(modbus_scheme_id)

    return c.String(http.StatusOK, "")
}

// ★保存Modbus方案相关地址
func modbus_address_save_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    modbus_address := Modbus_address{}

    json.Unmarshal(buf[0:n], &modbus_address)

    _MODBUS_MAIN.save_modbus_address(modbus_address)
    return c.String(http.StatusOK, "")
}

// ★获取Modbus的方案地址列表
func modbus_address_get_lst_restful(c echo.Context) error {
    modbus_scheme_id := c.Param("modbus_scheme_id")

    lst_modbus_address:= _MODBUS_MAIN.modbus_scheme_address_get_list(modbus_scheme_id)

    return c.JSON(http.StatusOK, lst_modbus_address)
}

// ★获取Modbus的方案地址对应的所有值
func modbus_address_and_value_lst_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    modbus_conf := Modbus_conf{}

    json.Unmarshal(buf[0:n], &modbus_conf)

    serial:=modbus_conf.Serial
    slaveId, _ := strconv.ParseUint(modbus_conf.SlaveId, 10, 8)
    baudrate, _ := strconv.Atoi(modbus_conf.BaudRate)
    parity, _ := strconv.Atoi(modbus_conf.Parity)
    dataBits, _ := strconv.Atoi(modbus_conf.DataBits)
    stopBits, _ := strconv.Atoi(modbus_conf.StopBits)
    scheme_id := modbus_conf.Scheme_id
    lst_modbus_address := _MODBUS_MAIN.modbus_scheme_address_and_value_get(serial,slaveId, baudrate, parity, dataBits, stopBits, scheme_id)

    return c.JSON(http.StatusOK, lst_modbus_address)
}

// ★设置Modbus的方案地址对应的所有值
func modbus_scheme_address_and_value_set_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    modbus_conf := Modbus_conf{}

    json.Unmarshal(buf[0:n], &modbus_conf)

    serial := modbus_conf.Serial
    slaveId, _ := strconv.ParseUint(modbus_conf.SlaveId, 10, 8)
    baudrate, _ := strconv.Atoi(modbus_conf.BaudRate)
    parity, _ := strconv.Atoi(modbus_conf.Parity)
    dataBits, _ := strconv.Atoi(modbus_conf.DataBits)
    stopBits, _ := strconv.Atoi(modbus_conf.StopBits)
    scheme_id := modbus_conf.Scheme_id
    _MODBUS_MAIN.modbus_scheme_address_and_value_set(serial, slaveId, baudrate, parity, dataBits, stopBits, scheme_id, modbus_conf.Modbus_address_and_value)

    return c.String(http.StatusOK, "")
}

// ★获取Modbus的地址
func modbus_address_get_restful(c echo.Context) error {
    modbus_scheme_id := c.Param("modbus_address_id")

    modbus_shceme_address := _MODBUS_MAIN.modbus_scheme_address_get(modbus_scheme_id)

    return c.JSON(http.StatusOK, modbus_shceme_address)
}

// ★删除Modbus方案地址
func modbus_address_delete_restful(c echo.Context) error {
    modbus_scheme_id := c.Param("modbus_address_id")

    _MODBUS_MAIN.modbus_scheme_address_delete(modbus_scheme_id)

    return c.String(http.StatusOK, "")
}

// ★一次性查询Modbus数据
func get_modbus_once_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    modbus_conf := Modbus_conf{}

    json.Unmarshal(buf[0:n], &modbus_conf)

    slaveId, _ := strconv.ParseUint(modbus_conf.SlaveId, 10, 8)
    baudrate, _ := strconv.Atoi(modbus_conf.BaudRate)
    parity, _ := strconv.Atoi(modbus_conf.Parity)
    data_type, _ := strconv.Atoi(modbus_conf.Modbus_Data_type)
    address, _ := strconv.ParseUint(modbus_conf.Address, 10, 32)

    f := _MODBUS_MAIN.get_modbus_once(modbus_conf.Serial, byte(slaveId), baudrate, parity, uint16(address-1), modbus_conf.Function_value, data_type)

    return c.JSON(http.StatusOK, f)
}

// ★一次性设置Modbus数据
func set_modbus_once_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    modbus_conf := Modbus_conf{}

    json.Unmarshal(buf[0:n], &modbus_conf)

    slaveId, _ := strconv.ParseUint(modbus_conf.SlaveId, 10, 8)
    baudrate, _ := strconv.Atoi(modbus_conf.BaudRate)
    parity, _ := strconv.Atoi(modbus_conf.Parity)
    data_type, _ := strconv.Atoi(modbus_conf.Modbus_Data_type)
    address, _ := strconv.ParseUint(modbus_conf.Address, 10, 32)

    _MODBUS_MAIN.set_modbus_once(modbus_conf.Serial, byte(slaveId), baudrate, parity, uint16(address-1), modbus_conf.Function_value, data_type, modbus_conf.Modbus_value)

    return c.String(http.StatusOK, "")
}

// ★保存Modbus设备
func modbus_device_save_restful(c echo.Context) error {
    buf := make([]byte, 1024)
    // ！！！这里是关键点！！！读取Request Payload的内容
    n, _ := c.Request().Body.Read(buf)

    ahathallan_modbus_Device := Ahathallan_modbus_Device{}

    json.Unmarshal(buf[0:n], &ahathallan_modbus_Device)

    ahathallan_modbus_device_main := Ahathallan_modbus_device_main{}
    ahathallan_modbus_device_main.init()

    ahathallan_modbus_device_main.save_Ahathallan_modbus_Device(ahathallan_modbus_Device)
    return c.String(http.StatusOK, "")
}

// ★获取Modbus设备列表
func modbus_device_get_lst_restful(c echo.Context) error {
    ahathallan_modbus_device_main := Ahathallan_modbus_device_main{}
    ahathallan_modbus_device_main.init()

    iot_equiment_lst:= ahathallan_modbus_device_main.get_lst_Ahathallan_modbus_Device()

    return c.JSON(http.StatusOK, iot_equiment_lst)
}

// ★获取Modbus设备
func modbus_device_get_restful(c echo.Context) error {
    modbus_device_id := c.Param("modbus_device_id")

    ahathallan_modbus_device_main := Ahathallan_modbus_device_main{}
    ahathallan_modbus_device_main.init()

    iot_equiment:= ahathallan_modbus_device_main.get_Ahathallan_modbus_Device(modbus_device_id)

    return c.JSON(http.StatusOK, iot_equiment)
}

// ★获取Modbus设备对应的值
func get_Ahathallan_modbus_Device_data_restful(c echo.Context) error {
    modbus_device_id := c.Param("modbus_device_id")

    ahathallan_modbus_device_main := Ahathallan_modbus_device_main{}
    ahathallan_modbus_device_main.init()

    iot_equiment := ahathallan_modbus_device_main.get_Ahathallan_modbus_Device_data(modbus_device_id)

    return c.String(http.StatusOK, iot_equiment)
}

//网络访问进程
func main_web() {
    //  处理异常的函数
    defer func() {
        fmt.Println("开始处理Web异常")
        // 获取异常信息
        if err := recover(); err != nil {
            //  输出异常信息
            errStr := fmt.Sprint("Web Error:", err)
            fmt.Println(errStr)
            mylooger.Output(2, errStr)
        } else {
            fmt.Println("无Web异常中止")
        }
        fmt.Println("结束Web异常处理")
    }()

    // Echo instance
    e := echo.New()

    e.Static("/static", "assets")

    // ★ 判断是否连接
    e.GET("/modbus_islink", modbus_islink_restful) // Routes

    // ★ 连接
    e.GET("/modbus_link", modbus_link_restful) // Routes

    // ★ 删除Modbus方案
    e.GET("/modbus_scheme_delete/:modbus_scheme_id", modbus_scheme_delete_restful) // Routes

    // ★ 获取Modbus方案
    e.GET("/modbus_scheme_get/:modbus_scheme_id", modbus_scheme_get_restful) // Routes

    // ★ 获取Modbus方案列表
    e.GET("/modbus_scheme_get_lst", modbus_scheme_get_lst_restful) // Routes

    // ★ 保存Modbus方案
    e.POST("/modbus_scheme_save", modbus_scheme_save_restful) // Routes

    // ★ 获取Modbus地址
    e.GET("/modbus_address_get/:modbus_address_id", modbus_address_get_restful) // Routes

    // ★ 删除Modbus地址
    e.GET("/modbus_address_delete/:modbus_address_id", modbus_address_delete_restful) // Routes

    // ★ 获取Modbus地址列表
    e.GET("/modbus_address_get_lst/:modbus_scheme_id", modbus_address_get_lst_restful) // Routes

    // ★ 获取Modbus地址列表
    e.POST("/modbus_address_and_value_lst", modbus_address_and_value_lst_restful) // Routes

    // ★ 设置Modbus地址列表和值
    e.POST("/modbus_scheme_address_and_value_set", modbus_scheme_address_and_value_set_restful) // Routes

    // ★ 保存Modbus方案相关地址
    e.POST("/modbus_address_save", modbus_address_save_restful) // Routes

    // ★ 断开
    e.GET("/modbus_dis_link", modbus_dis_link_restful) // Routes

    // ★ 一次性获取modbus值
    e.POST("/modbus_once_get", get_modbus_once_restful) // Routes

    // ★ 一次性设置modbus值
    e.POST("/modbus_once_set", set_modbus_once_restful) // Routes

    // ★ 保存Modbus设备
    e.POST("/modbus_device_save", modbus_device_save_restful) // Routes

    // ★ 获取Modbus设备列表
    e.GET("/modbus_device_get_lst", modbus_device_get_lst_restful) // Routes

    // ★ 获取Modbus设备
    e.GET("/modbus_device_get/:modbus_device_id", modbus_device_get_restful) // Routes

    // ★ 获取Modbus设备
    e.GET("/modbus_Device_data/:modbus_device_id", get_Ahathallan_modbus_Device_data_restful) // Routes

    //angulars
    //static_path := "./A003_Atom_Edge/web/"
    static_path := IOT_CONF.Web_folder

    e.Static("/js", static_path+"js")

    e.Static("/views", static_path+"views")

    e.Static("/css", static_path+"css")

    e.Static("/fonts", static_path+"fonts")

    e.Static("/img", static_path+"img")

    e.Static("/map", static_path+"map")

    e.Static("/login_static", static_path+"login2_files")

    e.Static("/data", static_path+"data")

    e.Static("/font-awesome", static_path+"font-awesome")

    e.File("/web", static_path+"index.html")

    e.File("/login", static_path+"login.html")

    e.File("/login_starcraft", static_path+"login_starcraft.html")

    // Start server
    //e.Logger.Fatal(e.Start(":1323"))

    socket := ":" + IOT_CONF.Web_socket
    e.Start(socket)
}



                            //■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■//

                            //↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓//

                                                                //Modbus操作方法

                            //↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑//

                            //■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■//

type Modbus_conf struct {
    Serial           string
    SlaveId          string
    BaudRate         string
    DataBits         string
    Parity           string
    StopBits         string
    Address          string
    Function_value   string
    Modbus_Data_type string
    Scheme_id        string

    Modbus_value             interface{}
    Modbus_address_and_value []Modbus_address_and_value_set
}

var _MODBUS_MAIN=modbus_main{}

type modbus_main struct {
    Serial      string
    SlaveId     byte
    BaudRate    int
    DataBits    int
    Parity      string
    StopBits    int
    Time_second int

    Is_link bool

    Handler *modbus.RTUClientHandler
}

type Modbus_scheme struct {
    Scheme_code string `json:"Scheme_code"`
    Device_Name string `json:"Device_Name"`
    Device_info string `json:"Device_info"`
    Device_type string `json:"Device_type"`
    SlaveId     string `json:"SlaveId"`
    BaudRate    string `json:"BaudRate"`
    DataBits    string `json:"DataBits"`
    Parity      string `json:"Parity"`
    StopBits    string `json:"StopBits"`
}

type Modbus_address struct {
    Address_scheme_code   string `json:"Address_scheme_code"`
    Address_No            string `json:"Address_No"`
    Address_Name          string `json:"Address_Name"`
    Address_Json          string `json:"Address_Json"`
    Address               string `json:"Address"`
    Function_value        string `json:"Function_value"`
    Function_Name         string `json:"Function_Name"`
    Modbus_Data_type      string `json:"Modbus_Data_type"`
    Modbus_Data_type_Name string `json:"Modbus_Data_type_Name"`
    Option_context        string `json:"Option_context"`
    UOM                   string `json:"UOM"`
}

type Modbus_address_and_value struct {
    Address_scheme_code   string `json:"Address_scheme_code"`
    Address_No            string `json:"Address_No"`
    Address_Name          string `json:"Address_Name"`
    Address_Json          string `json:"Address_Json"`
    Address               string `json:"Address"`
    Function_value        string `json:"Function_value"`
    Function_Name         string `json:"Function_Name"`
    Modbus_Data_type      string `json:"Modbus_Data_type"`
    Modbus_Data_type_Name string `json:"Modbus_Data_type_Name"`
    Option_context        string `json:"Option_context"`
    UOM                   string `json:"UOM"`

    Modbus_value interface{}
}

type Modbus_address_and_value_set struct {
    Address_scheme_code string `json:"Address_scheme_code"`
    Address_No          string `json:"Address_No"`
    Address_Name        string `json:"Address_Name"`
    Address             string `json:"Address"`
    Value               interface{}
}

//串口号
type Serial_Type int32

const (
    COM1      Serial_Type = 1
    COM2      Serial_Type = 2
    COM3      Serial_Type = 3
    COM4      Serial_Type = 4
    COM5      Serial_Type = 5
    COM6      Serial_Type = 6
    COM7      Serial_Type = 7
    COM8      Serial_Type = 8
)

//数据类型
type Data_type int

const (
    Unsigned_Integer      Data_type = 0
    Integer      Data_type = 1
    Double_Presion      Data_type = 2
    IEEE_Floating_Point      Data_type = 3
    IEEE_Reserved_World      Data_type = 4
)

func (p Data_type) String() string {
    switch (p) {
    case Unsigned_Integer:
        return "Unsigned_Integer"
    case Integer:
        return "Integer"
    case Double_Presion:
        return "Double_Presion"
    case IEEE_Floating_Point:
        return "IEEE_Floating_Point"
    case IEEE_Reserved_World:
        return "IEEE_Reserved_World"
    default:
        return "IEEE_Floating_Point"
    }
}

//获取配置文件
func (modbus_Main *modbus_main) init() {
    modbus_Main.Is_link = false
}

//字节数组转十六进制字符串
func (modbus_Main *modbus_main) bytes_hex_to_string(b []byte) string {
    return hex.EncodeToString(b)
}

//浮点f8转字节,float32转为字节数组,Modbus的4位float
func (modbus_Main *modbus_main) Float32bytes(float float32) []byte {
    bits := math.Float32bits(float)
    bytes := make([]byte, 4)
    binary.BigEndian.PutUint32(bytes, bits)
    return bytes
}

//Modbus 03H 指令，读取单个保持寄存器
///01 03 03 EC 00 02 05 BA  //  01=从机地址 03=功能码03位读 03 EC(1004 露点地址) 为寄存器地址 00 02 为寄存器数量
///01 03 04 41 80 E8 24 A1 FC  // 01=从机地址 03=功能码03位读 04 字节总数 41 80 E8 24 为数据（0,10000011,00000001110100000100100)
func (modbus_Main *modbus_main) Read_single_float(client modbus.Client,modbus_address uint16) float32 {
    results2, err := client.ReadHoldingRegisters(modbus_address, 2)

    if err == nil {
        //log.Println("label = ", i)
    } else {
        mylooger.Output(2, err.Error())
    }

    result_str := modbus_Main.bytes_hex_to_string(results2)

    n, _ := strconv.ParseUint(result_str, 16, 32)
    f := math.Float32frombits(uint32(n))
    fmt.Println("读取03指令,数据:%f", f)

    return f
}

//Modbus 10H 指令，读取多个保持寄存器WriteMultipleRegisters
//01 10 07d8 0002 04 c2a00000 e4ff //01=从机地址 10=功能码读取多个保持寄存器 07 D8(2008 露点地址) 为寄存器地址 00 02 为寄存器数量 04 字节数  c2a00000(float值)
func (modbus_Main *modbus_main) Write_multiple_float(client modbus.Client,modbus_address uint16,modbus_value float32) {
    //FLoat转字节
    s3 := modbus_Main.Float32bytes(modbus_value)
    //显示成16进制字符串
    fmt.Println("10H写入的字节16进制:%s", modbus_Main.bytes_hex_to_string(s3))

    results, err := client.WriteMultipleRegisters(modbus_address, 2, s3)

    if err == nil {

    } else {
        fmt.Sprint(results)
    }
}

//连接Modbus
func (modbus_Main *modbus_main) link_modbus() {
    // 声明modbus连接
    modbus_Main.Handler = modbus.NewRTUClientHandler(modbus_Main.Serial)
    modbus_Main.Handler.BaudRate = modbus_Main.BaudRate
    modbus_Main.Handler.DataBits = modbus_Main.DataBits
    modbus_Main.Handler.Parity = modbus_Main.Parity
    modbus_Main.Handler.StopBits = modbus_Main.StopBits
    modbus_Main.Handler.Timeout = time_second * 2 * time.Second

    // 声明modbus站号
    modbus_Main.Handler.SlaveId = modbus_Main.SlaveId

    _ = modbus_Main.Handler.Connect()
    modbus_Main.Is_link = true

    defer func() {
        modbus_Main.Is_link = false
        modbus_Main.Handler.Close()
    }()

    client := modbus.NewClient(modbus_Main.Handler)

    for {
        if modbus_Main.Is_link {
            fmt.Println(client)
            time.Sleep(time_second * 10 * time.Second)
        }
    }

    modbus_Main.Is_link = false
    modbus_Main.Handler.Close()
}

//一次性获取串口的Modbus数据
func (modbus_Main *modbus_main) get_modbus_once(serial string,slaveId byte,baudRate int,parity int, address uint16,function_value string,data_type int) interface{} {
    // 声明modbus连接
    modbus_Main.Handler = modbus.NewRTUClientHandler(serial)
    modbus_Main.Handler.BaudRate = baudRate
    modbus_Main.Handler.DataBits = 8
    if parity == 0 {
        modbus_Main.Handler.Parity = "N"
    } else if parity == 1 {
        modbus_Main.Handler.Parity = "P"
    } else if parity == 2 {
        modbus_Main.Handler.Parity = "E"
    }
    modbus_Main.Handler.StopBits = 1
    modbus_Main.Handler.Timeout = time_second * 2 * time.Second

    // 声明modbus站号
    modbus_Main.Handler.SlaveId = slaveId

    _ = modbus_Main.Handler.Connect()

    defer func() {
        modbus_Main.Handler.Close()
    }()

    client := modbus.NewClient(modbus_Main.Handler)

    //如果数据类型为2
    if data_type == 2 {

    }

    var value interface{} = 0

    if function_value == "03H" && data_type == 3 {
        //读取03H，Float数据
        value = modbus_Main.Read_single_float(client, address)
        fmt.Println(fmt.Sprintf("主操作：读取03指令,数据:%f", value))
    } else {
        //读取03H，Float数据
        f := modbus_Main.Read_single_float(client, address)
        fmt.Println(fmt.Sprintf("主操作：读取03指令,数据:%f", f))
    }

    modbus_Main.Handler.Close()

    return value
}

//一次性获取串口的Modbus数据
func (modbus_Main *modbus_main) set_modbus_once(serial string,slaveId byte,baudRate int,parity int, address uint16,function_value string,data_type int,modbus_value interface{}) {
    // 声明modbus连接
    modbus_Main.Handler = modbus.NewRTUClientHandler(serial)
    modbus_Main.Handler.BaudRate = baudRate
    modbus_Main.Handler.DataBits = 8
    if parity == 0 {
        modbus_Main.Handler.Parity = "N"
    } else if parity == 1 {
        modbus_Main.Handler.Parity = "P"
    } else if parity == 2 {
        modbus_Main.Handler.Parity = "E"
    }
    modbus_Main.Handler.StopBits = 1
    modbus_Main.Handler.Timeout = time_second * 2 * time.Second

    // 声明modbus站号
    modbus_Main.Handler.SlaveId = slaveId

    _ = modbus_Main.Handler.Connect()

    defer func() {
        modbus_Main.Handler.Close()
    }()

    client := modbus.NewClient(modbus_Main.Handler)

    if function_value == "10H" && data_type == 3 {
        //10H，Float数据
        f_str, _ := strconv.ParseFloat(modbus_value.(string), 64)

        modbus_Main.Write_multiple_float(client, address, float32(f_str))
    } else {
        //10H，Float数据
        f_str, _ := strconv.ParseFloat(modbus_value.(string), 64)

        modbus_Main.Write_multiple_float(client, address, float32(f_str))
    }

    modbus_Main.Handler.Close()
}

//获取配置文件
func (modbus_Main *modbus_main) dis_link_modbus() {
    modbus_Main.Is_link = false
    modbus_Main.Handler.Close()
}

//获取配置文件
func (modbus_Main *modbus_main) is_value() bool {
    return modbus_Main.Is_link
}

//获取Modbus的方案列表
func (modbus_Main *modbus_main) modbus_scheme_get_lst() []*Modbus_scheme{
    web_bolt_opt := new_App_bolt_opt("modbus_scheme_block")

    search_db_lst :=web_bolt_opt.get_all_db(reflect.ValueOf(Modbus_scheme{}).Type())

    c_hurt_data_lst := []*Modbus_scheme{}
    for i := 0; i < len(search_db_lst); i++ {
        modbus_scheme:=search_db_lst[i].(*Modbus_scheme)
        c_hurt_data_lst = append(c_hurt_data_lst, modbus_scheme)
        //fmt.Println(search_db_lst[i].(*Com_hurt_data))
    }

    return c_hurt_data_lst
}

//获取Modbus的方案
func (modbus_Main *modbus_main) modbus_scheme_get(id string) Modbus_scheme {
    modbus_scheme := Modbus_scheme{}

    web_bolt_opt := new_App_bolt_opt("modbus_scheme_block")

    web_bolt_opt.get_key_db(id, &modbus_scheme)

    return modbus_scheme
}

//获取Modbus的方案
func (modbus_Main *modbus_main) modbus_scheme_delete(id string) {
    web_bolt_opt := new_App_bolt_opt("modbus_scheme_block")

    web_bolt_opt.delete_db(id)
}

//保存Modbus的方案信息
func (modbus_Main *modbus_main) save_modbus_schedual(modbus_scheme Modbus_scheme) {
    web_bolt_opt := new_App_bolt_opt("modbus_scheme_block")

    if modbus_scheme.Scheme_code == "" {
        u1 := uuid.Must(uuid.NewV4())

        modbus_scheme.Scheme_code = u1.String()
    }

    web_bolt_opt.save_db(modbus_scheme.Scheme_code, modbus_scheme)
}

//获取Modbus的方案地址列表,查询条件是方案ID
func (modbus_Main *modbus_main) modbus_scheme_address_get_list(id string) []*Modbus_address {
    web_bolt_opt := new_App_bolt_opt("modbus_address_block")

    search_db_lst := web_bolt_opt.get_all_db(reflect.ValueOf(Modbus_address{}).Type())

    c_hurt_data_lst := []*Modbus_address{}
    for i := 0; i < len(search_db_lst); i++ {
        modbus_address := search_db_lst[i].(*Modbus_address)
        if modbus_address.Address_scheme_code == id {
            c_hurt_data_lst = append(c_hurt_data_lst, modbus_address)
        }

        if modbus_address.Function_value=="03H"{
            modbus_address.Function_Name="03H读保持寄存器"
        }else if modbus_address.Function_value=="10H"{
            modbus_address.Function_Name="10H写多个保持寄存器"
        }

        if modbus_address.Modbus_Data_type=="3"{
            modbus_address.Modbus_Data_type_Name="float"
        }else if modbus_address.Modbus_Data_type=="5"{
            modbus_address.Modbus_Data_type_Name="选择"
        }else if modbus_address.Modbus_Data_type=="0"{
            modbus_address.Modbus_Data_type_Name="int"
        }
        //fmt.Println(search_db_lst[i].(*Com_hurt_data))
    }

    return c_hurt_data_lst
}

//获取Modbus的地址
func (modbus_Main *modbus_main) modbus_scheme_address_get(id string) Modbus_address {
    modbus_address := Modbus_address{}

    web_bolt_opt := new_App_bolt_opt("modbus_address_block")

    web_bolt_opt.get_key_db(id, &modbus_address)

    return modbus_address
}

//获取Modbus的地址
func (modbus_Main *modbus_main) modbus_scheme_address_delete(id string) {
    web_bolt_opt := new_App_bolt_opt("modbus_address_block")

    web_bolt_opt.delete_db(id)
}

//★★获取Modbus的方案地址列表对应的所有值,查询条件是方案ID
func (modbus_Main *modbus_main) modbus_scheme_address_and_value_get(serial string,slaveId uint64,baudrate int,parity int, dataBits int,stopBits int,scheme_id string) []*Modbus_address_and_value {
    web_bolt_opt := new_App_bolt_opt("modbus_address_block")

    search_db_lst := web_bolt_opt.get_all_db(reflect.ValueOf(Modbus_address_and_value{}).Type())

    c_hurt_data_lst := []*Modbus_address_and_value{}
    for i := 0; i < len(search_db_lst); i++ {
        modbus_address := search_db_lst[i].(*Modbus_address_and_value)
        if modbus_address.Address_scheme_code == scheme_id {
            c_hurt_data_lst = append(c_hurt_data_lst, modbus_address)

            function_name := "03H"

            //读写都会加入读的值
            if modbus_address.Function_value == "03H" {
                modbus_address.Function_Name = "03H读保持寄存器"
            } else if modbus_address.Function_value == "10H" {
                modbus_address.Function_Name = "10H写多个保持寄存器"
                function_name = "03H"
            }

            if modbus_address.Modbus_Data_type == "3" {
                modbus_address.Modbus_Data_type_Name = "float"
            } else if modbus_address.Modbus_Data_type == "5" {
                modbus_address.Modbus_Data_type_Name = "选择"
            } else if modbus_address.Modbus_Data_type == "0" {
                modbus_address.Modbus_Data_type_Name = "int"
            }
            //fmt.Println(search_db_lst[i].(*Com_hurt_data))

            //取得对应的值
            address, _ := strconv.ParseUint(modbus_address.Address, 10, 32)
            data_type, _ := strconv.Atoi(modbus_address.Modbus_Data_type)

            f := modbus_Main.get_modbus_once(serial, byte(slaveId), baudrate, parity, uint16(address)-1, function_name, data_type)

            modbus_address.Modbus_value = f
        }
    }

    return c_hurt_data_lst
}

//★★获取Modbus的方案地址列表设置对应值,查询条件是方案ID
func (modbus_Main *modbus_main) modbus_scheme_address_and_value_set(serial string,slaveId uint64,baudrate int,parity int, dataBits int,stopBits int,scheme_id string,modbus_address_and_value []Modbus_address_and_value_set){
    address_get_list:= _MODBUS_MAIN.modbus_scheme_address_get_list(scheme_id)

    for i := 0; i < len(address_get_list); i++ {
        for j := 0; j < len(modbus_address_and_value); j++ {
            if address_get_list[i].Address_No==modbus_address_and_value[j].Address_No {
                address, _ := strconv.ParseUint(address_get_list[i].Address, 10, 32)
                data_type, _ := strconv.Atoi(address_get_list[i].Modbus_Data_type)

                modbus_Main.set_modbus_once(serial, byte(slaveId), baudrate, parity, uint16(address)-1, address_get_list[i].Function_value, data_type, modbus_address_and_value[j].Value)
            }
        }
    }
}

//保存Modbus的地址信息
func (modbus_Main *modbus_main) save_modbus_address(modbus_address Modbus_address) {
    web_bolt_opt := new_App_bolt_opt("modbus_address_block")

    if modbus_address.Address_No=="" {
        u1 := uuid.Must(uuid.NewV4())

        modbus_address.Address_No = u1.String()
    }
    web_bolt_opt.save_db(modbus_address.Address_No, modbus_address)
}

                            //■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■//

                            //↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓//

                                                               //设备管理

                            //↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑//

                            //■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■//

//设备元数据定义------------↓↓↓↓------------------------设备元数据定义------------↓↓↓↓------------------------设备元数据定义------------↓↓↓↓------------------------
type Ahathallan_modbus_device_main struct {
    Db_name string
}

func (ahathallan_modbus_device_main *Ahathallan_modbus_device_main) init() {
    ahathallan_modbus_device_main.Db_name="device_block"
}

//Modbus设备
type Ahathallan_modbus_Device struct {
    Dev_No      string `json:"Dev_No"`
    Dev_Name    string `json:"Dev_Name"`
    Scheme_code string `json:"Scheme_code"`
    Scheme_name string `json:"Scheme_name"`
    Serial      string `json:"Serial"`
}

func (ahathallan_modbus_device_main *Ahathallan_modbus_device_main) get_lst_Ahathallan_modbus_Device() []*Ahathallan_modbus_Device {

    var iot_equiment_lst []*Ahathallan_modbus_Device

    //从数据库取
    web_bolt_opt := new_App_bolt_opt_query(ahathallan_modbus_device_main.Db_name)

    search_db_lst := web_bolt_opt.get_all_db(reflect.ValueOf(Ahathallan_modbus_Device{}).Type())

    for i := 0; i < len(search_db_lst); i++ {
        ahathallan_modbus_Device := search_db_lst[i].(*Ahathallan_modbus_Device)

        iot_equiment_lst = append(iot_equiment_lst, ahathallan_modbus_Device)
    }

    return iot_equiment_lst
}

func (ahathallan_modbus_device_main *Ahathallan_modbus_device_main) get_Ahathallan_modbus_Device(dev_No string) Ahathallan_modbus_Device {
    iot_equiment := Ahathallan_modbus_Device{}

    //从数据库取
    web_bolt_opt := new_App_bolt_opt_query(ahathallan_modbus_device_main.Db_name)

    web_bolt_opt.get_key_db(dev_No, &iot_equiment)

    return iot_equiment
}

func(ahathallan_modbus_device_main *Ahathallan_modbus_device_main) save_Ahathallan_modbus_Device (ahathallan_modbus_Device Ahathallan_modbus_Device) {

    bolt_opt := new_App_bolt_opt(ahathallan_modbus_device_main.Db_name)

    bolt_opt.save_db(ahathallan_modbus_Device.Dev_No, ahathallan_modbus_Device)
}

func(ahathallan_modbus_device_main *Ahathallan_modbus_device_main) delete_Ahathallan_modbus_Device (dev_name string) {

    bolt_opt := new_App_bolt_opt_query(ahathallan_modbus_device_main.Db_name)

    bolt_opt.delete_db(dev_name)
}

//★★获取Modbus的设备对应值,查询条件是设备ID
func (ahathallan_modbus_device_main *Ahathallan_modbus_device_main) get_Ahathallan_modbus_Device_data(dev_No string) string {
    //获取设备
    Ahathallan_modbus_Device := ahathallan_modbus_device_main.get_Ahathallan_modbus_Device(dev_No)

    //获取设备对应的Modbus方案
    modbus_Main := modbus_main{}
    modbus_Main.init()

    modbus_shceme := _MODBUS_MAIN.modbus_scheme_get(Ahathallan_modbus_Device.Scheme_code)

    serial := Ahathallan_modbus_Device.Serial
    slaveId, _ := strconv.ParseUint(modbus_shceme.SlaveId, 10, 8)
    baudrate, _ := strconv.Atoi(modbus_shceme.BaudRate)
    parity, _ := strconv.Atoi(modbus_shceme.Parity)
    dataBits, _ := strconv.Atoi(modbus_shceme.DataBits)
    stopBits, _ := strconv.Atoi(modbus_shceme.StopBits)
    scheme_id := modbus_shceme.Scheme_code

    //取得Modbus方案对应的值
    lst_modbus_address := _MODBUS_MAIN.modbus_scheme_address_and_value_get(serial, slaveId, baudrate, parity, dataBits, stopBits, scheme_id)

    json_str := `{`

    for i := 0; i < len(lst_modbus_address); i++ {
        t := lst_modbus_address[i]
        json_str += `"` + t.Address_Json + `":`
        if i != len(lst_modbus_address)-1 {
            json_str += `"` + fmt.Sprint(t.Modbus_value) + `",`
        } else {
            json_str += `"` + fmt.Sprint(t.Modbus_value) + `"}`
        }
    }

    return json_str
}

//设备元数据定义-------------↑↑↑↑-----------------------设备元数据定义-------------↑↑↑↑-----------------------设备元数据定义-------------↑↑↑↑-----------------------

func open_web(){
    time.Sleep(3 * time.Second)
    // 有GUI调用
    exec.Command(`cmd`, `/c`, `start`, `http://localhost:13836/web`).Start()
}

func main() {
    //  处理异常的函数
    defer func() {
        fmt.Println("开始处理异常")
        // 获取异常信息
        if err := recover(); err != nil {
            //  输出异常信息
            errStr := fmt.Sprint("MAIN_ERR:", err)
            fmt.Println(errStr)
            mylooger.Output(2, errStr)
        } else {
            fmt.Println("无异常中止")
        }
        fmt.Println("结束异常处理")
    }()

    //SetFlags 重新设置输出选项
    mylooger = LogInfo();
    mylooger.SetFlags(log.Ldate | log.Ltime)

    //读取配置文件
    IOT_CONF = get_conf()
    time_second = time.Duration(IOT_CONF.Time_second)

    _MODBUS_MAIN.init()

    //go _MODBUS_MAIN.link_modbus()

    go open_web()

    main_web()
}
/*↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑主函数入口↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑主函数入口↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
#############################################################*/