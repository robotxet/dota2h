(function ($) {

var btn = document.getElementById("btnSubmit");
var info = document.getElementById("info")
var spinner = document.getElementById("spinner")

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
                $.ajax({
                    url: "/load_image",
                    type: "POST",
                    data: e.target.result,
                    contentType:"image/" + extension + "; base64",
                    success: function (filename) {
                        localFilename = filename
                    },
                });
                document.getElementById("loaded_img").style.visibility = "visible"; 
                $('#loaded_img').attr('src', e.target.result);
                btn.disabled = false;

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
    spinner.style.visibility = "visible";
    $.ajax({
        url: "/process_tf",
        type: "POST",
        data: localFilename,
        contentType: "text/plain",
        success: function (result) {
            result.sort(compare)
            var totalPercentage = 0;
            for (var i = 0, len = result.length; i < len; i++) {
                totalPercentage += result[i]["Rating"]
            }

            // TODO for each sum rating and make percentile + scale avatars based on persentile
            console.log(result)
            $("#heroname").html(
                '<div class="column-left">' + '<div>' + result[0]["Hero"] + ' : ' + (result[0]["Rating"] / totalPercentage).toFixed(2) + '</div>' + '</div>' +
                '<div class="column-center">' + '<div>' + result[1]["Hero"] + ' : ' + (result[1]["Rating"] / totalPercentage).toFixed(2) + '</div>' + '</div>' +
                '<div class="column-right">' + '<div>' + result[2]["Hero"] + ' : ' + (result[2]["Rating"] / totalPercentage).toFixed(2) + '</div>'
                )

            $("#avatar").html(
                '<div class="column-left"><div class="image-cropper"><img id="avatar_img" src="data:image/png;base64,' + result[0]["ImgData"] + '"/></div></div>' +
                '<div class="column-center"><div class="image-cropper-second"><img id="avatar_img" src="data:image/png;base64,' + result[1]["ImgData"] + '"/></div></div>' +
                '<div class="column-right"><div class="image-cropper-third"><img id="avatar_img" src="data:image/png;base64,' + result[2]["ImgData"] + '"/></div></div>'
                )
            $("#history").html(
                '<div class="column-left-history">' + result[0]["History"] + '</div>' +
                '<div class="column-center-history">' + result[1]["History"] + '</div>' +
                '<div class="column-right-history">' + result[2]["History"] + '</div>'
                )
        },
    });
    $(document).ajaxComplete(function() {
        spinner.style.visibility = "hidden";
        info.style.visibility = "visible";
    });

    e.preventDefault();

}

$("#btnSubmit").change(function(){
    calcTf(this);
});

}(jQuery));