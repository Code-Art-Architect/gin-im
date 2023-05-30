const userIdName = "userId"
const userInfoName = "userInfo"

function userId(id) {
  if (typeof id == 'undefined') {
    const r = sessionStorage.getItem(userIdName)
    if (!r) {
      return 0
    } else {
      return parseInt(r)
    }
  } else {
    sessionStorage.setItem(userIdName, id)
  }
}

function userInfo(o) {
  if (typeof o == 'undefined') {
    const r = sessionStorage.getItem(userInfoName)
    if (!!r) {
      return JSON.parse(r)
    } else {
      return null
    }
  } else {
    sessionStorage.setItem(userInfoName, JSON.stringify(o))
  }
}