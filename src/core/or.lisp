(macro (xs)
       ((lambda (y ys)
          (cond
	   y (quote t)
	   (list ys) (recur (car ys) (cdr ys))
	   (quote t) ())) (car xs) (cdr xs)))
