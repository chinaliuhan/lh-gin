const URL_LOGIN = "/user/login"
const URL_REGISTER = "/user/register"
const URL_GA_GENERATE = "/user/generateGa"
const URL_GA_QRCODE = "/user/generateGaQrcode"
const URL_GA_BIND = "/user/gaBind"

const CODE_SUCCESS = 0
const CODE_FAILED = 0

class common {
}

common.generateGAQrcode = function (id, name, text) {
    let qrcode = new QRCode(document.getElementById(id), {
        text: "otpauth://totp/" + name + "?secret=" + text,
        width: 128,
        height: 128,
        colorDark: "#000000",
        colorLight: "#ffffff",
        correctLevel: QRCode.CorrectLevel.H
    });

    // qrcode.clear(); // clear the code.
    // qrcode.makeCode("http://www.baidu.com"); // make another code.
}

common.post2Payload = function (path, params, callback) {
    $.ajax({
            url: path,
            type: "POST",     //规定请求的类型（GET 或 POST）。
            data: params,     //规定要发送到服务器的数据
            dataType: "json", //预期的服务器响应的数据类型
            async: true,      //布尔值，表示请求是否异步处理。默认是 true。
            contentType: "application/json;charset=utf-8", //发送数据到服务器时所使用的内容类型。默认是："application/x-www-form-urlencoded"。
            // "contentType": "application/x-www-form-urlencoded;charset=utf-8",
            timeout: 5000, //设置本地的请求超时时间（以毫秒计）。
            beforeSend: function (xhr) {       //发送请求前运行的函数。

            },
            complete: function (xhr, status, error) {     //请求完成时运行的函数（在请求成功或失败之后均调用，即在 success 和 error 函数之后）。
            },
            error: function (xhr, status, error) {
                console.error("请求失败", xhr, status, error)
            },
            success: function (xhr, status, error) {
                callback(xhr, status, error)
            }
        }
    )
}

common.post2Form = function (path, params, callback) {
    $.ajax({
            url: path,
            type: "POST",     //规定请求的类型（GET 或 POST）。
            data: params,     //规定要发送到服务器的数据
            dataType: "POST", //预期的服务器响应的数据类型
            async: true,      //布尔值，表示请求是否异步处理。默认是 true。
            // contentType: "application/json;charset=utf-8", //发送数据到服务器时所使用的内容类型。默认是："application/x-www-form-urlencoded"。
            "contentType": "application/x-www-form-urlencoded;charset=utf-8",
            timeout: 5000, //设置本地的请求超时时间（以毫秒计）。
            beforeSend: function (xhr) {       //发送请求前运行的函数。

            },
            complete: function (xhr, status, error) {     //请求完成时运行的函数（在请求成功或失败之后均调用，即在 success 和 error 函数之后）。
                console.log(xhr)
                console.log(status)
                console.log(error)
            },
            error: function (xhr, status, error) {

            },
            success: function (xhr, status, error) {

            }
        }
    )
}

common.getJson = function (path, callback) {

    $.get(path, callback, "json")
}

