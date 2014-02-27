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
package org.textmapper.tool.parser.ast;

import org.textmapper.tool.parser.TMTree.TextSource;

public class TmaRhsAssignment extends TmaNode implements ITmaRhsPart {

	private final TmaIdentifier id;
	private final boolean addition;
	private final ITmaRhsPart inner;

	public TmaRhsAssignment(TmaIdentifier id, boolean addition, ITmaRhsPart inner, TextSource source, int offset, int endoffset) {
		super(source, offset, endoffset);
		this.id = id;
		this.addition = addition;
		this.inner = inner;
	}

	public TmaIdentifier getId() {
		return id;
	}

	public boolean getAddition() {
		return addition;
	}

	public ITmaRhsPart getInner() {
		return inner;
	}

	public void accept(TmaVisitor v) {
		if (!v.visit(this)) {
			return;
		}
		if (id != null) {
			id.accept(v);
		}
		if (inner != null) {
			inner.accept(v);
		}
	}
}