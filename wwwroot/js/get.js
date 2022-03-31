function get(url) {
    return new Promise(function (resolve, reject) {
        let request = new XMLHttpRequest();
        request.open('Get', url, true);
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
        request.send();
    })
}
