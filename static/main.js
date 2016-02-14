$(document).ready(function() {
  $('.remove-button').on('click', function(event) {
    var target = $(event.target);
    target.closest("li").remove();
    var repoId = target.attr('id');
    var starred = target.data('starred');

    var data = { "id": repoId, "starred": starred };
    $.post($SCRIPT_ROOT + '/ignore', data);
  });

  $('.remove-all').on('click', function(event) {
    if (confirm('X all on this page?')) {
      $('.remove-button').trigger('click');
    }
  });
});
