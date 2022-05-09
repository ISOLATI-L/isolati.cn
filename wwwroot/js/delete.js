function deleteReq(url, data, contentType) {
    return new Promise(function (resolve, reject) {
        let request = new XMLHttpRequest();
        request.open('DELETE', url, true);
        if (contentType != "") {
            request.setRequestHeader("Content-type", contentType);
        }
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
