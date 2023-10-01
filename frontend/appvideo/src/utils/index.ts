export function getFormatTime(dateTime, flag) {
  if (dateTime != null) {
    var time = parseInt(String(dateTime * 1000))
    var date = new Date(time)
    var YY = date.getFullYear()
    var MM = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1
    var DD = date.getDate() < 10 ? '0' + date.getDate() : date.getDate()
    if (flag) {
      var hh = date.getHours() < 10 ? '0' + date.getHours() : date.getHours()
      var mm = date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()
      var ss = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds()
      return YY + '-' + MM + '-' + DD + ' ' + hh + ':' + mm + ':' + ss
    } else {
      return YY + '-' + MM + '-' + DD
    }
  } else {
    return ''
  }
}
