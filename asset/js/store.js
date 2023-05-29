function userId(id) {
  if (typeof id == 'undefined') {
    const r = sessionStorage.getItem('userid')
    if (!r) {
      return 0
    } else {
      return parseInt(r)
    }
  } else {
    sessionStorage.setItem('userid', id)
  }
}

function userInfo(o) {
  if (typeof o == 'undefined') {
    const r = sessionStorage.getItem('userinfo')
    if (!!r) {
      return JSON.parse(r)
    } else {
      return null
    }
  } else {
    sessionStorage.setItem('userinfo', JSON.stringify(o))
  }
}

const url = location.href
const isOpen = url.indexOf('/login') > -1 || url.indexOf('/register') > -1
if (!userId() && !isOpen) {
  // location.href = "login.shtml";
}

export {userId, userInfo}