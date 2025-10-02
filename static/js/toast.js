var toastElList = [].slice.call(document.querySelectorAll('.toast[autoshow]'))

var toastList = toastElList.map(function (toastEl) {
	toastEl.removeAttribute('autoshow');
var toast = new bootstrap.Toast(toastEl)
	toast.show()
return toast
})

toastElList.forEach(function(toastEl) {
	toastEl.addEventListener('hidden.bs.toast', function() {
		this.remove()
	})
})

