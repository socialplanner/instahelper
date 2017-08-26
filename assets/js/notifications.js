$(document).ready(function() {
  // Remove the number that displays the number of notifications on click.
  // Dont remove timeout
  $("#notifications").on("click", function() {
    setTimeout($(".notification").remove, 50);
  });
});
