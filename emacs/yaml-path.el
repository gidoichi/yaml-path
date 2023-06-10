;;; yaml-path --- a yaml element path that user/cursor indicated

;;; Commentary:

;; yaml-path.el explain the path of an element in any yaml file.
;; yaml-path.el needs exectable yaml-path file.

;;; Code:

(defcustom yaml-path-bin
  "yaml-path"
  "Path to yaml-path binary."
  :group 'yaml-path
  :type 'string
  :safe 'stringp)

(defcustom yaml-path-output-format
  "bosh"
  "Output format of yaml-path, i.e. argument of --format of yaml-path command."
  :group 'yaml-path
  :type 'string
  :safe 'stringp)

;;;###autoload
(defun yaml-path-at-point()
  "Get and display yaml path at point."
  (interactive)
  (let* ((path (yaml-path-get-path-at-point)))
    (kill-new path)
    (message "%s" path))
  )

(defun yaml-path-get-path-at-point(&optional pline pcol)
  "Get and return yaml path at point (pline as PLINE, pcol as PCOL)."
  (let ((result "???")
        (line (if pline pline (number-to-string (line-number-at-pos))))
        (col  (if pcol  pcol  (number-to-string (current-column))))
        (outbuf (get-buffer-create "*yaml-path-result*")))
    (cond ((zerop (progn
                    (with-current-buffer outbuf (erase-buffer))
                    (call-process-region
                     (point-min) (point-max) yaml-path-bin nil outbuf nil
                     "--line" line "--col" col "--format" yaml-path-output-format)))
           (with-current-buffer outbuf
             (setq result (replace-regexp-in-string "\n+" "" (buffer-string)))))
          ((zerop (progn
                    (with-current-buffer outbuf (erase-buffer))
                    (call-process-region
                     (point-min) (point-max) yaml-path-bin nil outbuf nil
                     "--line" line "--format" yaml-path-output-format)))
           (with-current-buffer outbuf
             (setq result (replace-regexp-in-string "\n+" "" (buffer-string))))))
    (kill-buffer outbuf)
    result
    ))

;;;###autoload
(defun yaml-path-which-func()
  "Set yaml-path to which-function's hook."
  (add-hook 'which-func-functions 'yaml-path-get-path-at-point t t)
  )

;; --------------------------------------------------------------------------- ;

;;;###autoload
(put 'yaml-path-bin 'safe-local-variable 'stringp)

(provide 'yaml-path)

;; Local Variables:
;; ispell-local-dictionary: "american"
;; End:

;;; yaml-path.el ends here
