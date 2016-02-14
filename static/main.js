$(document).ready(function() {
  $('.remove-button').on('click', function(event) {
    $(event.target).closest("li").remove();
    var repoId = $(event.target).attr('id');

    var data = { "id": repoId };
    $.post($SCRIPT_ROOT + '/ignore', data);
  });

  $('.remove-all').on('click', function(event) {
    $('.remove-button:not(:last)').trigger('click');
  });
});
