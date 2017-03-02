% CC-TEST-REPORTER-SUM-COVERAGE(1) User Manuals
% Code Climate <hello@codeclimate.com>
% February 2017

# PROLOG

This is a sub-command of **cc-test-reporter**(1).

# SYNOPSIS

**cc-test-reporter-sum-coverage** [--output=\<path>] FILE [FILE, ...]

# DESCRIPTION

Combine (sum) multiple pre-formatted coverage payloads into one.

# OPTIONS

## -o, --output *PATH*

Output to *PATH*. If *-* is given, content will be written to *stdout*. Defaults
to *coverage/codeclimate.json*.

## FILE [FILE, ...]

Input files to combine. These are expected to be pre-formatted coverage
payloads. Passing a single file will return it unprocessed.

# ALGORITHM

## SUM-ABILITY

The following must be true for payloads to be sum-able. If these conditions are
not met, an error will be returned:

1. The value for *git.head* is equal across all payloads
1. All *source_files[].coverage* arrays for the same *name* are the same length

## METADATA

Some keys will not differ between partial payloads, or their differences do not
matter to our system. They are taken from the first payload given on the
command-line.

    ci_service.*
    environment.*
    git.*
    run_at

## SOURCE FILES

All *source_files* values of the same *name* are combined by a reduce or fold
operation with the following binary operator, expressed in pseudo-code:

    combine_source_files a b = {
      name = a[name],
      blob_id: a[blob_id],
      coverage: combine_coverage a b,
      line_counts: {
        missed: count is-zero coverage,
        total: count non-null coverage,
        covered: total - missed,
      },
      covered_percent: line_counts.covered / line_counts.total * 100,
      covered_strength: TODO,
    }

    combine_coverage a b = for index in [0 .. length a[coverage]] {
      a_value = a[coverage][index]
      b_value = b[coverage][index]

      if a_value && b_value
        a_value + b_value
      else
        a_value || b_value
      end
    }

The results of these by-name folds are then concatenated.

## COVERAGE TOTALS

Given the summed *source_files* elements, top-level coverage attributes are
computed as followed:

    line_counts.missed:  sum source_files[].line_counts.missed
    line_counts.total:   sum source_files[].line_counts.total
    line_counts.covered: line_counts.total - line_counts.missed

    covered_percent: line_counts.covered / line_counts.total * 100
    covered_strength: TODO

# ENVIRONMENT VARIABLES

None

# SEE ALSO

**cc-test-reporter-format-coverage**(1).
