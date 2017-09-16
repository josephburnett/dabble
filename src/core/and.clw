(cons (cons (quote and)
            (cons (macro (x xs)
                         ((lambda (y ys)
                                  (cond
                                    y (cond
                                        (eq () ys) 1
                                        1 (recur (car ys) (cdr ys)))
                                    1 ())) x xs)) ())) ())
