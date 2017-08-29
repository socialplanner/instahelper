$(document).ready(function() {
  // Remove the number that displays the number of notifications on click.
  // Dont remove timeout
  $("#notifications").on("click", function() {
    setTimeout(function() {
      $(".notification").remove();
    }, 50);
  });
});

function showNotification(message, vertical, horizontal, type) {
  // type = ['','info','success','warning','danger'];
  // vertical = ['top', 'bottom']
  // horizontal = ['center', 'left', 'right']
  $.notify(
    {
      icon: "notifications",
      message: message
    },
    {
      type: type,
      timer: 4000,
      placement: {
        from: vertical,
        align: horizontal
      }
    }
  );
}
