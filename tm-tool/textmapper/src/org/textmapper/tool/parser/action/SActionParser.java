/**
 * Copyright 2002-2013 Evgeny Gryaznov
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.textmapper.tool.parser.action;

import java.io.IOException;
import java.text.MessageFormat;
import org.textmapper.tool.parser.action.SActionLexer.ErrorReporter;
import org.textmapper.tool.parser.action.SActionLexer.LapgSymbol;
import org.textmapper.tool.parser.action.SActionLexer.Lexems;

public class SActionParser {

	public static class ParseException extends Exception {
		private static final long serialVersionUID = 1L;

		public ParseException() {
		}
	}

	private final ErrorReporter reporter;

	public SActionParser(ErrorReporter reporter) {
		this.reporter = reporter;
	}

	private static final boolean DEBUG_SYNTAX = false;
	private static final int[] tmAction = SActionLexer.unpack_int(11,
		"\uffff\uffff\ufffd\uffff\ufff7\uffff\ufff1\uffff\4\0\uffff\uffff\uffff\uffff\3\0" +
		"\2\0\5\0\ufffe\uffff");

	private static final short[] tmLalr = SActionLexer.unpack_short(18,
		"\1\uffff\3\1\uffff\ufffe\1\uffff\3\1\uffff\ufffe\1\uffff\3\0\uffff\ufffe");

	private static final short[] lapg_sym_goto = SActionLexer.unpack_short(9,
		"\0\0\4\4\6\7\11\14\16");

	private static final short[] lapg_sym_from = SActionLexer.unpack_short(14,
		"\0\1\2\3\5\6\0\1\2\1\2\3\1\2");

	private static final short[] lapg_sym_to = SActionLexer.unpack_short(14,
		"\1\2\2\2\10\11\12\3\3\4\4\7\5\6");

	private static final short[] lapg_rlen = SActionLexer.unpack_short(6,
		"\1\0\3\2\1\3");

	private static final short[] lapg_rlex = SActionLexer.unpack_short(6,
		"\7\7\4\5\5\6");

	protected static final String[] lapg_syms = new String[] {
		"eoi",
		"'{'",
		"_skip",
		"'}'",
		"javaaction",
		"command_tokens",
		"command_token",
		"command_tokensopt",
	};

	public interface Tokens extends Lexems {
		// non-terminals
		public static final int javaaction = 4;
		public static final int command_tokens = 5;
		public static final int command_token = 6;
		public static final int command_tokensopt = 7;
	}

	/**
	 * -3-n   Lookahead (state id)
	 * -2     Error
	 * -1     Shift
	 * 0..n   Reduce (rule index)
	 */
	protected static int tmAction(int state, int symbol) {
		int p;
		if (tmAction[state] < -2) {
			if (symbol == Lexems.Unavailable_) {
				return -3 - state;
			}
			for (p = -tmAction[state] - 3; tmLalr[p] >= 0; p += 2) {
				if (tmLalr[p] == symbol) {
					break;
				}
			}
			return tmLalr[p + 1];
		}
		return tmAction[state];
	}

	protected static int tmGoto(int state, int symbol) {
		int min = lapg_sym_goto[symbol], max = lapg_sym_goto[symbol + 1] - 1;
		int i, e;

		while (min <= max) {
			e = (min + max) >> 1;
			i = lapg_sym_from[e];
			if (i == state) {
				return lapg_sym_to[e];
			} else if (i < state) {
				min = e + 1;
			} else {
				max = e - 1;
			}
		}
		return -1;
	}

	protected int lapg_head;
	protected LapgSymbol[] lapg_m;
	protected LapgSymbol lapg_n;
	protected SActionLexer lapg_lexer;

	public Object parse(SActionLexer lexer) throws IOException, ParseException {

		lapg_lexer = lexer;
		lapg_m = new LapgSymbol[1024];
		lapg_head = 0;

		lapg_m[0] = new LapgSymbol();
		lapg_m[0].state = 0;
		lapg_n = lapg_lexer.next();

		while (lapg_m[lapg_head].state != 10) {
			int lapg_i = tmAction(lapg_m[lapg_head].state, lapg_n == null ? Lexems.Unavailable_ : lapg_n.symbol);
			if (lapg_i <= -3 && lapg_n == null) {
				lapg_n = lapg_lexer.next();
				lapg_i = tmAction(lapg_m[lapg_head].state, lapg_n.symbol);
			}

			if (lapg_i >= 0) {
				reduce(lapg_i);
			} else if (lapg_i == -1) {
				shift();
			}

			if (lapg_i == -2 || lapg_m[lapg_head].state == -1) {
				break;
			}
		}

		if (lapg_m[lapg_head].state != 10) {
			reporter.error(lapg_n == null ? lapg_lexer.getOffset() : lapg_n.offset, lapg_n == null ? lapg_lexer.getLine() : lapg_lexer.getTokenLine(),
						MessageFormat.format("syntax error before line {0}",
								lapg_lexer.getTokenLine()));
			throw new ParseException();
		}
		return lapg_m[lapg_head].value;
	}

	protected void shift() throws IOException {
		if (lapg_n == null) {
			lapg_n = lapg_lexer.next();
		}
		lapg_m[++lapg_head] = lapg_n;
		lapg_m[lapg_head].state = tmGoto(lapg_m[lapg_head - 1].state, lapg_n.symbol);
		if (DEBUG_SYNTAX) {
			System.out.println(MessageFormat.format("shift: {0} ({1})", lapg_syms[lapg_n.symbol], lapg_lexer.current()));
		}
		if (lapg_m[lapg_head].state != -1 && lapg_n.symbol != 0) {
			lapg_n = null;
		}
	}

	protected void reduce(int rule) {
		LapgSymbol lapg_gg = new LapgSymbol();
		lapg_gg.value = (lapg_rlen[rule] != 0) ? lapg_m[lapg_head + 1 - lapg_rlen[rule]].value : null;
		lapg_gg.symbol = lapg_rlex[rule];
		lapg_gg.state = 0;
		if (DEBUG_SYNTAX) {
			System.out.println("reduce to " + lapg_syms[lapg_rlex[rule]]);
		}
		LapgSymbol startsym = (lapg_rlen[rule] != 0) ? lapg_m[lapg_head + 1 - lapg_rlen[rule]] : lapg_n;
		lapg_gg.offset = startsym == null ? lapg_lexer.getOffset() : startsym.offset;
		applyRule(lapg_gg, rule, lapg_rlen[rule]);
		for (int e = lapg_rlen[rule]; e > 0; e--) {
			lapg_m[lapg_head--] = null;
		}
		lapg_m[++lapg_head] = lapg_gg;
		lapg_m[lapg_head].state = tmGoto(lapg_m[lapg_head - 1].state, lapg_gg.symbol);
	}

	@SuppressWarnings("unchecked")
	protected void applyRule(LapgSymbol lapg_gg, int rule, int ruleLength) {
	}
}
