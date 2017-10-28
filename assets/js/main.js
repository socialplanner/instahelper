String.prototype.format = function() {
  var args = [].slice.call(arguments);
  return this.replace(/(\{\d+\})/g, function(a) {
    return args[+a.substr(1, a.length - 2) || 0];
  });
};

$(function() {
  $("form").submit(function(e) {
    var data = $(this).serialize();
    var url = $(this).attr("action");
    var method = $(this).attr("method");

    e.preventDefault();

    var button = $(this).find("button");
    button.prop("disabled", true);

    if (method.toLowerCase().startsWith("get")) {
      $(".spinner").show();
      $.get(url, data, function(d) {
        $(".spinner").hide();

        button.prop("disabled", false);

        checkData(d, d.status);
      }).fail(function() {
        $(".spinner").hide();
      });
    } else if (method.toLowerCase().startsWith("post")) {
      $(".spinner").show();
      $.post(url, data, function(d) {
        $(".spinner").hide();

        button.prop("disabled", false);

        checkData(d, d.status);
      }).fail(function(d) {
        $(".spinner").hide();

        button.prop("disabled", false);

        checkData(d.responseText, d.status);
      });
    }
  });

  StartWebSocket = function() {
    if (window.WebSocket) {
      var socket = new WebSocket(getWebSocketURL());
      window.onbeforeunload = function(event) {
        socket.close();
      };

      socket.onclose = function() {
        setTimeout(function() {
          StartWebSocket();
        }, 5000);
      };

      socket.onmessage = function(e) {
        var data = JSON.parse(e.data);
        showNotification(data.text, "top", "center", "info", data.link, 2000);
      };
    }
  };

  StartWebSocket();

  // Remove the number that displays the number of notifications on click.
  // Dont remove timeout
  $("#notifications").on("click", function() {
    setTimeout(function() {
      $(".notification").remove();
    }, 50);
  });

  $("#clear-all").on("click", function() {
    var i = $(this);

    deleteNotifications(function() {
      $(".notification").remove();

      i.siblings().remove();
    });
  });
});

function showNotification(message, vertical, horizontal, type, url, timer) {
  // type = ['','info','success','warning','danger'];
  // vertical = ['top', 'bottom']
  // horizontal = ['center', 'left', 'right']
  $.notify(
    {
      icon: "notifications",
      message: message,
      url: url
    },
    {
      type: type,
      timer: timer || 4000,
      placement: {
        from: vertical,
        align: horizontal
      }
    }
  );
}

// Displays error if res.status != 200 else display success
function checkData(d, status) {
  if (status != 200 && status !== undefined) {
    showNotification("Error: " + d, "top", "left", "danger");
  } else {
    showNotification(d, "top", "left", "success");
  }
}

function deleteNotifications(callback) {
  $.ajax({
    url: "/api/notifications",
    type: "DELETE",
    success: function(d) {
      callback(true, d);
    },
    error: function(d) {
      callback(false, d.responseText);
    }
  });
}

// Deletes the account from the database
function deleteAccount(username, callback) {
  $.ajax({
    url: "/api/accounts/" + username,
    type: "DELETE",
    success: function(d) {
      callback(true, d);
    },
    error: function(d) {
      callback(false, d.responseText);
    }
  });
}

// Returns a websocket url of the form url:port/ws
function getWebSocketURL() {
  var loc = window.location,
    url;
  if (loc.protocol === "https:") {
    url = "wss:";
  } else {
    url = "ws:";
  }
  url += "//" + loc.host;
  url += "/ws";

  return url;
}

$(".update-button").on("click", function() {
  var button = $(this);
  var version = button.attr("version");
  var description = button.attr("description");
  var url = "/api/update/to/" + version;

  swal({
    title: "Are you sure you want to update to {0}?".format(version),
    html: "<b>Description: </b> " + description,
    type: "warning",
    showCancelButton: true,
    confirmButtonColor: "#3085d6",
    cancelButtonColor: "#d33",
    confirmButtonText: "Yes, upgrade!"
  }).then(function() {
    $.ajax({
      url: url,
      type: "POST",
      success: function(d) {
        showNotification(
          "Awesome! Updated to " +
            version +
            "! All you need to do is restart the current running app.",
          "top",
          "left",
          "success"
        );
      },

      error: function(d) {
        showNotification("Error: " + d.responseText, "top", "left", "danger");
      }
    });
  });
});

$(".delete-account").on("click", function() {
  var button = $(this);
  var username = button.attr("username");
  swal({
    title: "Are you sure you would like to remove {0}?".format(username),
    text: "You won't be able to revert this!",
    type: "warning",
    showCancelButton: true,
    confirmButtonColor: "#3085d6",
    cancelButtonColor: "#d33",
    confirmButtonText: "Yes, delete it!"
  }).then(function() {
    deleteAccount(username, function(ok, res) {
      if (ok) {
        // Remove the row from the table
        button.closest("tr").remove();
        swal("Deleted!", "{0} has been deleted.".format(username), "success");
      } else {
        swal("Error!", res.format(username), "error");
      }
    });
  });
});

$("#jobs-submit").on("click", function() {
  var btn = $(this);

  // All forms with the class jobs-settings
  var data = $(".jobs-settings").serialize();

  btn.prop("disabled", true);
  $(".spinner").show();
  $.ajax({
    url: btn.attr("url"),
    type: "POST",
    success: function(d) {
      showNotification("Updated!", "top", "left", "success");
      btn.prop("disabled", false);
      $(".spinner").hide();
    },

    error: function(d) {
      showNotification("Error: " + d.responseText, "top", "left", "danger");
      btn.prop("disabled", false);
      $(".spinner").hide();
    }
  });
});
