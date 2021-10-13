(macro ((xs))
       (if (eq () xs) t
	 (label x (car xs)
		'(if (eq () `x) ()
		   `(apply recur (cdr xs))))))
