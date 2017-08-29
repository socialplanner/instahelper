$(function() {
  function checkData(d, status) {
    if (status != 200 && status !== undefined) {
      showNotification("Error: " + d, "top", "left", "danger");
    } else {
      showNotification(d, "top", "left", "success");
    }
  }

  $("form").submit(function(e) {
    var data = $(this).serialize();
    var url = $(this).attr("action");
    var method = $(this).attr("method");
    e.preventDefault();
    var button = $(this).find("button");
    button.prop("disabled", true);
    if (method.toLowerCase() == "get") {
      $.get(url, data, function(d) {
        button.prop("disabled", false);
        checkData(d, d.status);
      });
    } else if (method.toLowerCase() == "post") {
      $.post(url, data, function(d) {
        button.prop("disabled", false);
        checkData(d, d.status);
      }).fail(function(d) {
        button.prop("disabled", false);
        checkData(d.responseText, d.status);
      });
    }
  });
});
