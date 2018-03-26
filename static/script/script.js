var countdown;

$(document).ready(function() {
  countdown = $(".countdown").TimeCircles({
    "start": false,
    "total_duration": 5,
    "count_past_zero": false,
    "animation": "smooth",
    "bg_width": 1.2,
    "fg_width": 0.1,
    "circle_bg_color": "#60686F",
    "time": {
        "Days": {
            "show": false
        },
        "Hours": {
            "show": false
        },
        "Minutes": {
            "show": false
        },
        "Seconds": {
            "color": "#FF9999",
            "show": true
        }
    }
  }).addListener(countdownComplete);

});

function startCountdown() {
  $('.countdownwrapper').fadeIn();
  countdown.restart()
}

function countdownComplete(unit, value, total) {
  if (total<=0) { 
    $.ajax({
      type: 'post',
      url: '/photo',
      dataType: 'json',
      data: null,
      success: function(data) {
        if (data !== null) {
          $('.thumbnails').prepend('<img onmouseover="getElementById(\'preview\').src=this.alt" src="' + data.lastImage.thumb + '" alt="' + data.lastImage.image + '">');
        }
        $('.countdownwrapper').fadeOut();
      }
    });
  }
}