function post(url, data) {
    return new Promise(function (resolve, reject) {
        let request = new XMLHttpRequest();
        request.open('POST', url, true);
        request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        request.onload = function () {
            if (request.readyState == 4) {
                resolve(request);
            } else {
                reject(request);
            }
        };
        request.onerror = function () {
            reject(Error("Network Error"));
        };
        request.send(data);
    })
}
