(cons (cons (quote and)
            (cons (macro (x xs)
                         ((lambda (y ys)
                                  (cond
                                    y (cond
                                        (eq () ys) (quote t)
                                        (quote t) (recur (car ys) (cdr ys)))
                                    (quote t) ())) x xs)) ())) ())
