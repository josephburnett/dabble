(macro ((xs))
       '(if (eq () `(car xs))
	    ()
	  (if (eq () `(cdr xs))
	      t
	    (recur `(cdr xs)))))
