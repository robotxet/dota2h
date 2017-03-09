(function ($) { 

var fileTypes = ['jpg', 'jpeg', 'png'];
var localFilename;

function readURL(input) {
    if (input.files && input.files[0]) {
        var extension = input.files[0].name.split('.').pop().toLowerCase();
        isSuccess = fileTypes.indexOf(extension) > -1;
        if (isSuccess) {
            var file  = input.files[0];
            var reader = new FileReader();
          
            reader.onload = function (e) {
                // $('#loaded_img').attr('src', e.target.result);
            
                $.ajax({
                    url: "/load_image",
                    type: "POST",
                    data: e.target.result,
                    contentType:"image/" + extension + "; base64",
                    success: function (filename) {
                        localFilename = filename
                    },
                });
            }
            reader.readAsDataURL(file);
        } else {
            alert('wrong file type')
        }
    }
}
    
$("#file").change(function(){
    readURL(this);
});

var form = document.getElementById("calcTflow");

document.getElementById("calcTflow").addEventListener("submit", function (e) {
  calcTf(e);
});

function calcTf(e) {
    $.ajax({
        url: "/process_tf",
        type: "POST",
        data: localFilename,
        contentType: "text/plain",
        success: function (result) {
            $("#result").val(atob(result["TfData"]))
            $("#avatar").html('<div class="image-cropper"><img id="avatar_img" src="data:image/png;base64,' + result["ImgData"] + '"/></div>')
            $("#history").html('<textarea readonly rows=10 cols=80 id="result">' + result["History"] + '</textarea>')
        },
    });

    e.preventDefault();

}

$("#btnSubmit").change(function(){
    calcTf(this);
});

}(jQuery));