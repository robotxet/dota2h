(function ($) { 

var fileTypes = ['jpg', 'jpeg', 'png'];
var localFilename;

function readURL(input) {
    if (input.files && input.files[0]) {
        var extension = input.files[0].name.split('.').pop().toLowerCase();
        isSuccess = fileTypes.indexOf(extension) > -1;
        console.log(extension)
        if (isSuccess) {
            var file  = input.files[0];
            var reader = new FileReader();
          
            reader.onload = function (e) {
                $('#loaded_img').attr('src', e.target.result);
            
                $.ajax({
                    url: "/load_image",
                    type: "POST",
                    data: e.target.result,
                    contentType:"image/" + extension + "; base64",
                    success: function (filename) {
                        localFilename = filename
                        console.log(filename);
                    },
                });
            }
            reader.readAsDataURL(file);
        } else {
            alert('wrong file type')
        }
    }
}
    
$("#imgInp").change(function(){
    readURL(this);
});

var form = document.getElementById("calcTflow");

document.getElementById("calcTflow").addEventListener("submit", function (e) {
  calcTf(e);
});

function calcTf(e) {
    console.log("here")
    console.log(localFilename)
    $.ajax({
        url: "/process_tf",
        type: "POST",
        data: localFilename,
        contentType: "text/plain",
        success: function (result) {
            $("#result").val(result)
        },
    });

    e.preventDefault();

}

$("#btnSubmit").change(function(){
    calcTf(this);
});

}(jQuery));