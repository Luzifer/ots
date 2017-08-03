securePassword = null

createSecret = () ->
  secret = $('#formCreateSecret').find('textarea').val()

  if $('#extra').prop 'checked'
    securePassword = Math.random().toString(36).substring(2)
    secret = GibberishAES.enc(secret, securePassword)

  $.ajax 'api/create',
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

newSecret = () ->
  location.href = location.href.split('#')[0]
  false

secretCreated = (data) ->
  secretHash = data.secret_id
  if securePassword != null
    secretHash = "#{secretHash}|#{securePassword}"
  url = "#{location.href.split('#')[0]}##{secretHash}"

  $('#panelNewSecret').hide()
  $('#panelSecretURL').show()
  $('#panelSecretURL').find('input').val url
  $('#panelSecretURL').find('input').focus()
  $('#panelSecretURL').find('input').select()

  securePassword = null

showData = (data) ->
  secret =  data.secret
  if securePassword != null
    secret = GibberishAES.dec(secret, securePassword)

  $('#panelNewSecret').hide()
  $('#panelSecretURL').hide()
  $('#notfound').hide()
  $('#somethingwrong').hide()
  $('#panelReadSecret').show()
  $('#panelReadSecret').find('textarea').val secret

somethingWrong = () ->
  $('#somethingwrong').show()


$ ->
  initBinds()
  hashLoad()
