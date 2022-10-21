$(document).ready(function () {
  $('#toc-toggle-open').on('click', function () {
    $('.col-md-2-5').addClass('opened')
    $('body').addClass('modal-opened')
  })

  $('#toc-toggle-close').on('click', function () {
    $('.col-md-2-5').removeClass('opened')
    $('body').removeClass('modal-opened')
  })

  $('#toc a').on('click', function (){
    $('.col-md-2-5').removeClass('opened')
    $('body').removeClass('modal-opened')
  })

  /* Collapsing toc */
  $('#toc ul ul').addClass('collapse')

  // Change initial icon for nav without children:
  $('#toc .nav-collapse-handler').each(function () {
    if ($(this).siblings('ul').length === 0) {
      $(this).find('.glyphicon').removeClass('glyphicon-triangle-bottom')
      $(this).find('.glyphicon').addClass('glyphicon-chevron-down')
      $(this).addClass('no-children')
    }
  })

  // Expand / collpase on click
  $('#toc .nav-collapse-handler').on('click', function() {
    toggleNav($(this))
  })

  $(docSidebarInitialExpand)
})

// Expand / collpase on click
function toggleNav(el) {
  if (el.hasClass('collapsed')) {
    if (!el.hasClass('no-children')) {
      el.removeClass('collapsed')
      el.siblings('ul').collapse('show')
    }
  } else {
    el.addClass('collapsed')
    el.siblings('ul').collapse('hide')
  }
}

const docSidebarInitialExpand = function () {
  const toc = $('#toc')
  const pathname = window.location.pathname
  const hash = window.location.hash
  toc.find('a[href="'+pathname+hash+'"]').each(function(i, nav) {
    $(nav).parents('ul').each(function(i, el) {
      $(el).collapse('show')
      $(el).siblings('span.nav-collapse-handler:not(.no-children)').removeClass('collapsed')
      $(el).siblings('span.nav-collapse-handler').addClass('active')
    })
    $(nav).siblings('span.nav-collapse-handler:not(.no-children)').removeClass('collapsed')
    $(nav).siblings('span.nav-collapse-handler').addClass('active')
    $(nav).siblings('ul').collapse('show')
  })
}
