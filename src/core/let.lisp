(macro (bindings form)
       (if (eq () bindings) form
	 (label b (car bindings)
		'(label `(car b) `(car (cdr b)) `form))))
