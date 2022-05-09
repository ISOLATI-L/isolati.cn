let dialogMask = document.querySelector(".dialogMask");
let deleteDialog = document.querySelector(".dialog");
let information = document.querySelector("#information");
let deleteID;

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

function deleteParagraph() {
    information.innerHTML = "删除中　";
    deleteReq("/admin/api/paragraph", String(deleteID), "").then(
        function (res) {
            console.log(res);
            if (res.status === 200) {
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

