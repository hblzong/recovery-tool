(function (doc, win) {
  var resizeEvt =
    "orientationchange" in window ? "orientationchange" : "resize";

  function calc() {
    var dEl = document.documentElement,
      clientWidth = dEl.clientWidth;
    var maxWidth = 640,
      minWidth = 320;

    if (clientWidth > maxWidth) {
      clientWidth = maxWidth;
    } else if (clientWidth < minWidth) {
      clientWidth = minWidth;
    }
    // if(clientWidth>=640){
    // 	clientWidth = 640;
    // }

    dEl.style.fontSize = 100 * (clientWidth / 750) + "px";
  }
  function calc1920() {
    var dEl = document.documentElement,
      clientWidth = dEl.clientWidth;
    var maxWidth = 1920,
      minWidth = 1024;

    if (clientWidth > maxWidth) {
      clientWidth = maxWidth;
    } else if (clientWidth < minWidth) {
      clientWidth = minWidth;
    }
    dEl.style.fontSize = 100 * (clientWidth / 1920) + "px";
  }
  var dEl = document.documentElement,
    clientWidth = dEl.clientWidth;
  if (clientWidth <= 960) {
    win.addEventListener(resizeEvt, calc);
    doc.addEventListener("DOMContentLoaded", calc);
  } else {
    win.addEventListener(resizeEvt, calc1920);
    doc.addEventListener("DOMContentLoaded", calc1920);
  }
  $(".info-mobile").on("click", ".info-ml", function () {
    $(this).parent().toggleClass("active").siblings(".dropdown").slideToggle();
    $(".info-mobile").css("border-radius", "10px 10px 0px 0px");
  });
})(document, window);
