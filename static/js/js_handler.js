(function ($) { 

function compare(a,b) {
  if (a.Rating > b.Rating)
    return -1;
  if (a.Rating < b.Rating)
    return 1;
  return 0;
}

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
            result.sort(compare)
            // TODO for each sum rating and make percentile + scale avatars based on persentile
            console.log(result)
            $("#heroname").html(
                '<div class="column-left">' + result[0]["Hero"] + ' : ' + result[0]["Rating"] + '</div>' +
                '<div class="column-left">' + result[1]["Hero"] + ' : ' + result[1]["Rating"] + '</div>' +
                '<div class="column-left">' + result[2]["Hero"] + ' : ' + result[2]["Rating"] + '</div>'
                )

            $("#avatar").html(
                '<div class="column-left"><div class="image-cropper"><img id="avatar_img" src="data:image/png;base64,' + result[0]["ImgData"] + '"/></div></div>' +
                '<div class="column-center"><div class="image-cropper"><img id="avatar_img" src="data:image/png;base64,' + result[1]["ImgData"] + '"/></div></div>' +
                '<div class="column-right"><div class="image-cropper"><img id="avatar_img" src="data:image/png;base64,' + result[2]["ImgData"] + '"/></div></div>'
                )
            var history = document.getElementById("history");
            history.textContent = result[0]["History"]
        },
    });

    e.preventDefault();

}

$("#btnSubmit").change(function(){
    calcTf(this);
});

}(jQuery));