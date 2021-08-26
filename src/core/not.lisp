((not (macro (x xs)
        (cond
	  (list xs) (abc "not requires one argument")
	  (quote t) (cond
	              x ()
		      (quote t) (quote t))))))
