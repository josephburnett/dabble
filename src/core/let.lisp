(cons (cons (quote let)
            (cons (macro (binding form xs)
                         (cond
			  (eq () (atom xs)) (error "Too many parameters to let.")
			  (eq () (quote binding)) (error "No bindings provided.")
			  (quote t) (label (car (car (quote binding)))
					   (car (cdr (quote binding)))
					   (cond
					    (eq () (cdr (quote binding))) form
					    (quote t) (recur (cdr (quote binding)) form))))) ())) ())
