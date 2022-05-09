let dialogMask = document.querySelector(".dialogMask");
let deleteDialog = document.querySelector(".dialog");
let information = document.querySelector("#information");
let uploadInformationFont = document.querySelector("#uploadInformationFont");
let deleteID;
let imgData = [];
let nImg = 0;
let loadedImg = 0;
let oParent = document.querySelector(".images");
let requesting = false;
let loading = false;

let imgFile = document.querySelector(".imgFile");

window.onload = function () {
    window.onscroll = loadImg;
    loadImg();
};

function loadImg() {
    if (!loading && !requesting && checkScroll()) {
        if (nImg == imgData.length) {
            requesting = true;
            get("/images/api/list?s=" + String(nImg) + "&n=10").then(
                function (res) {
                    if (res.status === 200) {
                        let data = JSON.parse(res.response);
                        if (data.length > 0) {
                            imgData = imgData.concat(data);
                            addElement();
                        } else {
                            window.onscroll = null;
                        }
                        requesting = false;
                    } else {
                        setTimeout(function () {
                            requesting = false;
                        }, 10000);
                    }
                },
                function (res) {
                    console.log(res);
                    setTimeout(function () {
                        requesting = false;
                    }, 10000);
                }
            );
        } else {
            addElement();
        }
    }
}

function addElement() {
    loading = true;
    let imgID = imgData[nImg].split(".")[0];
    let liElement = document.createElement("li");
    liElement.style.display = "none";
    liElement.id = "img" + imgID;
    oParent.appendChild(liElement);
    let aElement = document.createElement("a");
    aElement.className = "imageBox";
    aElement.target = "_blank";
    aElement.href = "/images/" + imgData[nImg];
    liElement.appendChild(aElement);
    let a2Element = document.createElement("a");
    a2Element.className = "default_link";
    a2Element.href = "javascript:showDeleteDialog(" + imgID + ");";
    a2Element.innerHTML = "删除";
    liElement.appendChild(a2Element);
    let imgElement = document.createElement("img");
    imgElement.className = "image";
    imgElement.src = "/images/" + imgData[nImg];
    imgElement.onload = function () {
        liElement.style.display = "inline-block";
        loading = false;
        loadedImg++;
        loadImg()
    }
    aElement.appendChild(imgElement);
    nImg++;
}

function checkScroll() {
    let imgBox = document.querySelectorAll('.images li');
    if (imgBox.length < 3) {
        return true;
    }
    let maxH = imgBox[imgBox.length - 3].offsetTop + imgBox[imgBox.length - 3].offsetHeight;
    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
    let documentHeight = document.documentElement.clientHeight;
    return maxH < scrollTop + documentHeight;
}

async function uploadImage() {
    let file = imgFile.files[0];
    if (!file) {
        return;
    }
    uploadInformationFont.innerHTML = "上传中";

    let fileData = {};
    fileData["suffix"] = "." + file.name.split(".").slice(-1)[0];
    fileData["data"] = await readFile(file);
    post("/admin/api/image", JSON.stringify(fileData), "application/json").then(
        function (res) {
            console.log(res);
            if (res.status === 200) {
                uploadInformationFont.innerHTML = "上传成功";
                loadImg();
            } else {
                uploadInformationFont.innerHTML = "上传失败";
            }
        },
        function (res) {
            console.log(res);
            uploadInformationFont.innerHTML = "上传失败";
        });
}

function readFile(file) {
    return new Promise(function (resolve) {
        let reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = function () {
            resolve(reader.result.split(",")[1]);
        };
    });
}

function cancelDelete() {
    dialogMask.style.display = "none";
    deleteDialog.style.display = "none";
}

function showDeleteDialog(id) {
    deleteID = id;
    information.innerHTML = "　　　　";
    dialogMask.style.display = "block";
    deleteDialog.style.display = "block";
}

function deleteImage() {
    information.innerHTML = "删除中　";
    deleteReq("/admin/api/image", String(deleteID), "").then(
        function (res) {
            console.log(res);
            if (res.status === 200) {
                let deletedElement = document.querySelector("#img" + deleteID);
                deletedElement.parentNode.removeChild(deletedElement);
                loadImg();
                information.innerHTML = "删除成功";
            } else {
                information.innerHTML = "删除失败";
            }
        },
        function (res) {
            console.log(res);
            information.innerHTML = "删除失败";
        });
}
