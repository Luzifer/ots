securePassword = null

createSecret = () ->
  secret = $('#formCreateSecret').find('textarea').val()

  securePassword = Math.random().toString(36).substring(2)
  secret = GibberishAES.enc(secret, securePassword)

  $.ajax 'api/create',
    method: "post"
    data:
      secret: secret
    dataType: "json"
    statusCode:
      201: secretCreated
      400: somethingWrong
      500: somethingWrong
      404: () ->
        # Mock for interface testing
        secretCreated
          secret_id: 'foobar'

  false

dataNotFound = () ->
  $('#notfound').show()

hashLoad = () ->
  hash = window.location.hash
  if hash.length == 0
    return

  $('#cardNewSecret').hide()
  $('#cardSecretURL').hide()
  $('#notfound').hide()
  $('#somethingwrong').hide()
  $('#cardReadSecretPre').show()

requestSecret = () ->
  hash = window.location.hash
  hash = decodeURIComponent(hash)

  parts = hash.split '|'
  if parts.length == 2
    hash = parts[0]
    securePassword = parts[1]

  id = hash.substring(1)
  $.ajax "api/get/#{id}",
    dataType: "json"
    statusCode:
      404: dataNotFound
      200: showData

initBinds = () ->
  $('#formCreateSecret').bind 'submit', createSecret
  $('#newSecret, .navbar-brand').bind 'click', newSecret
  $(window).bind 'hashchange', hashLoad
  $('#revealSecret').bind 'click', requestSecret
  bindResizeTextarea()

newSecret = () ->
  location.href = location.href.split('#')[0]
  false

bindResizeTextarea = () ->
  $('textarea').each(() ->
    text = this

    doResize = () =>
      #text.style.height = 'auto'
      text.style.height = (this.scrollHeight) + 'px'

    delayedResize = () =>
      window.setTimeout doResize, 0

    text.setAttribute('style', "height: #{this.scrollHeight}px; min-height: #{this.scrollHeight}px; overflow-y:hidden;")

    $(text)
      .on('change', doResize)
      .on('cut', delayedResize)
      .on('paste', delayedResize)
      .on('drop', delayedResize)
      .on('keydown', delayedResize)
  )


secretCreated = (data) ->
  secretHash = data.secret_id
  if securePassword != null
    secretHash = "#{secretHash}|#{securePassword}"
  url = "#{location.href.split('#')[0]}##{secretHash}"

  $('#cardNewSecret').hide()
  $('#cardReadSecretPre').hide()
  $('#cardSecretURL').show()
  $('#cardSecretURL').find('input').val url
  $('#cardSecretURL').find('input').focus()
  $('#cardSecretURL').find('input').select()

  securePassword = null

showData = (data) ->
  secret =  data.secret
  if securePassword != null
    secret = GibberishAES.dec(secret, securePassword)

  $('#cardNewSecret').hide()
  $('#cardSecretURL').hide()
  $('#notfound').hide()
  $('#somethingwrong').hide()
  $('#cardReadSecretPre').hide()
  $('#cardReadSecret').show()
  $('#cardReadSecret').find('textarea').val secret
  $('#cardReadSecret').find('textarea').trigger 'change'

somethingWrong = () ->
  $('#somethingwrong').show()


$ ->
  initBinds()
  hashLoad()
