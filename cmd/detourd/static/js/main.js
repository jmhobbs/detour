$(function () {
  var $mapping = $("#mapping"),
    $overlay = $(".overlay"),
    $formWindow = $(".overlay-window"),
    $loadingWindow = $(".overlay-loading-window"),
    $hostname = $("#hostname"),
    $address = $("#address");

  $.getJSON("/api/list", function (data) {
    refreshTable(data.hosts);

    $overlay.hide();
    $loadingWindow.hide();
  });

  function refreshTable (hosts) {
    $mapping.children().remove();
    $.each(hosts, function (hostname, address) {
      $mapping.append(tmpl("row-template", {hostname: hostname, address: address}));
    });
  }

  $(".icon-plus").on("click", function () {
    $hostname.val("");
    $address.val("");
    $overlay.show();
    $formWindow.show();
  });

  $("#form-cancel").on("click", function () {
    $overlay.hide();
    $formWindow.hide();
  });

  var current_hostname = null;

  $("#form-save").on("click", function () {
    var options = {host: $hostname.val(), ip: $address.val()};

    $formWindow.hide();
    $loadingWindow.show();

    // If we changed the hostname, we need to clear it out.
    if(current_hostname != null && options.host != current_hostname) {
      $.post("/api/unset", {host: current_hostname});
    }

    $.post("/api/set", options, function (data) {
      var mapping = JSON.parse(data);
      refreshTable(mapping.hosts);
      $loadingWindow.hide();
      $overlay.hide();
      current_hostname = null;
    });
  });

  $(document).on("click", ".icon-pencil", function () {
    var $tr = $(this).closest("tr"),
      hostname = $(this).closest("tr").find('.hostname').text(),
      address = $(this).closest("tr").find('.address').text();

    current_hostname = hostname;

    $hostname.val(hostname);
    $address.val(address);

    $overlay.show();
    $formWindow.show();
  });

  $(document).on("click", ".icon-cancel", function () {
    $overlay.show();
    $loadingWindow.show();
    $.post("/api/unset", { host: $(this).closest("tr").find(".hostname").text() }, function (data) {
      var mapping = JSON.parse(data);
      refreshTable(mapping.hosts);
      $loadingWindow.hide();
      $overlay.hide();
    });
  });
});
