let imgData = [];
const cols = 3;
let arrH = [-5, -5, -5];
let nImg = 0;
let loadedImg = 0;
let oParent = document.querySelector(".images");
let end = false;
let requesting = false;

window.onload = function () {
    window.onscroll = function () {
        loadImg();
    }
    loadImg();
};

function loadImg() {
    if (checkScroll()) {
        if (nImg == imgData.length) {
            if (!requesting) {
                requesting = true;
                get("/images/api/list?s=" + String(nImg) + "&n=10").then(
                    function (res) {
                        if (res.status === 200) {
                            let data = JSON.parse(res.response);
                            if (data.length > 0) {
                                imgData = imgData.concat(data);
                                addElement();
                            } else {
                                end = true;
                            }
                        } else {
                            end = true;
                        }
                        requesting = false;
                    },
                    function (res) {
                        console.log(res);
                        end = true;
                        requesting = false;
                    }
                );
            }
        } else {
            addElement();
        }
    }
}

function addElement() {
    let aElement = document.createElement("a");
    aElement.className = "imageBox";
    aElement.target = "_blank";
    aElement.style.display = "none";
    aElement.href = "/images/" + imgData[nImg];
    oParent.appendChild(aElement);
    let imgElement = document.createElement("img");
    imgElement.className = "image";
    imgElement.src = "/images/" + imgData[nImg];
    imgElement.onload = function () {
        waterfall(aElement, imgElement);
        loadImg()
    }
    aElement.appendChild(imgElement);
    nImg++;
}

function waterfall(aElement, imgElement) {
    let minH = Math.min.apply(null, arrH);
    let index = getIndex(arrH, minH);
    aElement.style.top = minH + 5 + "px";
    aElement.style.left = 285 * index + "px";
    aElement.style.display = "block";
    arrH[index] += imgElement.offsetHeight + 5;
    let maxH = Math.max.apply(null, arrH);
    oParent.style.height = maxH + "px";
    loadedImg++;
}

function getIndex(arr, val) {
    for (i in arr) {
        if (arr[i] == val) {
            return i;
        }
    }
}

function checkScroll() {
    if (end) {
        return false
    }
    let imgBox = document.querySelectorAll('.imageBox');
    if (imgBox.length == 0) {
        return true;
    }
    let minH = Math.min.apply(null, arrH);
    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
    let documentHeight = document.documentElement.clientHeight;
    return minH < scrollTop + documentHeight;
}
