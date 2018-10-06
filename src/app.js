let securePassword = null

// bindResizeTextarea attaches resize triggers to all available text areas
function bindResizeTextarea() {
  $('textarea').each((idx, text) => {
    let doResize = () => {
      text.style.height = text.scrollHeight + 'px'
    }

    let delayedResize = () => {
      window.setTimeout(doResize, 0)
    }

    text.setAttribute('style', "height: #{this.scrollHeight}px; min-height: #{this.scrollHeight}px; overflow-y:hidden;")

    $(text)
      .on('change', doResize)
      .on('cut', delayedResize)
      .on('paste', delayedResize)
      .on('drop', delayedResize)
      .on('keydown', delayedResize)
  })
}

// createSecret executes the secret creation after encrypting the secret
function createSecret() {
  let secret = $('#formCreateSecret').find('textarea').val()

  securePassword = Math.random().toString(36).substring(2)
  secret = GibberishAES.enc(secret, securePassword)

  $.ajax('api/create', {
    method: "post",
    data: {
      secret: secret,
    },
    dataType: "json",
    statusCode: {
      201: secretCreated,
      400: somethingWrong,
      500: somethingWrong,
      404: () => {
        // Mock for interface testing
        secretCreated({
          secret_id: 'foobar',
        })
      },
    },
  })

  return false
}

// dataNotFound displays the not-found error
function dataNotFound() {
  $('#notfound').show()
}

// hashLoad reacts on a changed window hash an starts the diplaying of the secret
function hashLoad() {
  let hash = window.location.hash
  if (hash.length === 0) return

  $('#cardNewSecret').hide()
  $('#cardSecretURL').hide()
  $('#notfound').hide()
  $('#somethingwrong').hide()
  $('#cardReadSecretPre').show()
}

// initBinds attaches functions to frontend elements
function initBinds() {
  $('#formCreateSecret').bind('submit', createSecret)
  $('#newSecret, .navbar-brand').bind('click', newSecret)
  $(window).bind('hashchange', hashLoad)
  $('#revealSecret').bind('click', requestSecret)
  bindResizeTextarea()
}

// newSecret removes the window hash and therefore returns to "new secret" mode
function newSecret() {
  location.href = location.href.split('#')[0]
}

// requestSecret requests the encrypted secret from the backend
function requestSecret() {
  let hash = window.location.hash
  hash = decodeURIComponent(hash)

  let parts = hash.split('|')
  if (parts.length === 2) {
    hash = parts[0]
    securePassword = parts[1]
  }

  let id = hash.substring(1)
  $.ajax(`api/get/${id}`, {
    dataType: "json",
    statusCode: {
      404: dataNotFound,
      200: showData,
    }
  })
}

// secretCreated generates the share URLs and displays them to the user
function secretCreated(data) {
  let secretHash = data.secret_id
  if (securePassword !== null) secretHash = `${secretHash}|${securePassword}`
  let url = `${location.href.split('#')[0]}#${secretHash}`

  $('#cardNewSecret').hide()
  $('#cardReadSecretPre').hide()
  $('#cardSecretURL').show()
  $('#cardSecretURL').find('input').val(url)
  $('#cardSecretURL').find('input').focus()
  $('#cardSecretURL').find('input').select()

  securePassword = null
}

// showData takes the backend answer, decrypts the secret and shows it
function showData(data) {
  let secret = data.secret
  if (securePassword !== null) secret = GibberishAES.dec(secret, securePassword)

  $('#cardNewSecret').hide()
  $('#cardSecretURL').hide()
  $('#notfound').hide()
  $('#somethingwrong').hide()
  $('#cardReadSecretPre').hide()
  $('#cardReadSecret').show()
  $('#cardReadSecret').find('textarea').val(secret)
  $('#cardReadSecret').find('textarea').trigger('change')
}

// somethingWrong shows the "something went wrong" screen
function somethingWrong() {
  $('#somethingwrong').show()
}

// Trigger initialization functions
$(() => {
  initBinds()
  hashLoad()
})
